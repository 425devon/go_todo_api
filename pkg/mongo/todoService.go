package mongo

import (
	"github.com/425devon/go_todo_api/pkg/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TodoService struct {
	collection *mgo.Collection
}

func todoIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func NewTodoService(session *Session, dbName string, collectionName string) *TodoService {
	collection := session.GetCollection(dbName, collectionName)
	collection.EnsureIndex(todoIndex())
	return &TodoService{collection}
}

func (s *TodoService) CreateList(list *models.TodoList) (_id string, err error) {
	list.ID = bson.NewObjectId()
	return list.ID.Hex(), s.collection.Insert(&list)
}

func (s *TodoService) CreateTask(list *models.TodoList, task *models.Task) (_id string, err error) {
	task.ID = bson.NewObjectId()
	task.Completed = false
	l, err := s.GetListByID(list.ID.Hex())
	if err != nil {
		return "", err
	}
	l.Tasks = append(l.Tasks, *task)
	return task.ID.Hex(), s.collection.UpdateId(l.ID, &l)
}

func (s *TodoService) GetListByID(_id string) (models.TodoList, error) {
	var todoList models.TodoList
	err := s.collection.FindId(bson.ObjectIdHex(_id)).One(&todoList)
	return todoList, err
}

func (s *TodoService) GetTaskByID(_id string) (models.Task, error) {
	var task models.Task
	query := bson.M{
		"tasks._id": bson.ObjectIdHex(_id),
	}
	err := s.collection.Find(query).One(&task)
	return task, err
}

//GetAllLists stil needs in test
func (s *TodoService) GetAllLists() ([]models.TodoList, error) {
	var lists []models.TodoList
	err := s.collection.Find(bson.M{}).All(&lists)
	return lists, err
}

func (s *TodoService) CompleteTask(_id string) (*models.Task, error) {
	var task models.Task
	query := bson.M{
		"tasks._id": bson.ObjectIdHex(_id),
	}

	update := bson.M{"$set": bson.M{"completed": true}}

	err := s.collection.Update(query, update)
	if err != nil {
		return nil, err
	}
	err = s.collection.Find(query).One(&task)

	return &task, err
}

func (s *TodoService) DeleteListByID(_id string) error {
	err := s.collection.RemoveId(bson.ObjectIdHex(_id))
	return err
}

//DeleteTaskByID still needs int test
func (s *TodoService) DeleteTaskByID(_id string) error {
	query := bson.M{
		"tasks._id": bson.ObjectIdHex(_id),
	}
	err := s.collection.Remove(query)
	return err
}
