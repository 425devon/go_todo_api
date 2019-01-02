package mongo

import (
	"errors"

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
	for _, tsk := range l.Tasks {
		if tsk.Name == task.Name {
			return "", errors.New("duplicate task")
		}
	}
	l.Tasks = append(l.Tasks, *task)
	return task.ID.Hex(), s.collection.UpdateId(l.ID, &l)
}

func (s *TodoService) GetListByID(_id string) (models.TodoList, error) {
	var todoList models.TodoList
	err := s.collection.FindId(bson.ObjectIdHex(_id)).One(&todoList)
	return todoList, err
}

func (s *TodoService) GetTaskByID(listID, taskID string) (*models.Task, error) {
	list, err := s.GetListByID(listID)
	for key, task := range list.Tasks {
		if task.ID == bson.ObjectIdHex(taskID) {
			tsk := list.Tasks[key]
			return &tsk, err
		}
	}
	return nil, err
}

//GetAllLists stil needs int test
func (s *TodoService) GetAllLists() ([]models.TodoList, error) {
	var lists []models.TodoList
	err := s.collection.Find(bson.M{}).All(&lists)
	return lists, err
}

func (s *TodoService) CompleteTask(listID, taskID string) (*models.Task, error) {
	list, err := s.GetListByID(listID)
	for key, task := range list.Tasks {
		if task.ID == bson.ObjectIdHex(taskID) {
			list.Tasks[key].Completed = true
			tsk := list.Tasks[key]
			err = s.collection.UpdateId(list.ID, list)
			return &tsk, err
		}
	}
	return nil, err
}

func (s *TodoService) DeleteListByID(_id string) error {
	return s.collection.RemoveId(bson.ObjectIdHex(_id))
}

//DeleteTaskByID still needs int test
func (s *TodoService) DeleteTaskByID(listID, taskID string) error {
	list, err := s.GetListByID(listID)
	if err != nil {
		return err
	}
	for key, task := range list.Tasks {
		if task.ID == bson.ObjectIdHex(taskID) {
			list.Tasks = append(list.Tasks[:key], list.Tasks[key+1:]...)
			return s.collection.UpdateId(list.ID, list)
		}
	}
	return errors.New("task not found")
}
