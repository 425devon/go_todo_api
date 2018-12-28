package mongo

import (
	"encoding/json"

	"github.com/425devon/go_todo_api/pkg/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TodoService struct {
	collection *mgo.Collection
}

func todoIndex() mgo.Index {
	return mgo.Index{
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

func (s *TodoService) CreateTask(list *models.TodoList, task *models.Task) (id string, err error) {
	task.ID = bson.NewObjectId()
	task.Completed = false
	l, err := s.GetListByID(list.ID.Hex())
	if err != nil {
		return "", err
	}
	l.Tasks = append(l.Tasks, *task)

	// pretty, _ := json.Marshal(l)
	// fmt.Println(string(pretty))
	return task.ID.Hex(), s.collection.UpdateId(l.ID, &l)
}

func (s *TodoService) GetListByID(_id string) (models.TodoList, error) {
	var todoList models.TodoList
	err := s.collection.FindId(bson.ObjectIdHex(_id)).One(&todoList)
	return todoList, err
}

func (s *TodoService) GetTaskByID(_id string) (models.Task, error) {
	var task models.Task
	taskID := bson.ObjectIdHex(_id)
	pipe := s.collection.Pipe([]bson.M{
		{"$match": bson.M{
			"tasks._id": taskID,
		}},
	})
	resp := bson.M{}
	err := pipe.One(&resp)
	if err != nil {
		return task, err
	}

	tasks := resp["tasks"].([]interface{})
	singleTask := tasks[0]

	tsk, err := json.Marshal(singleTask)
	if err != nil {
		return task, err
	}

	err = json.Unmarshal(tsk, &task)
	if err != nil {
		return task, err
	}

	return task, err
}

// func (s *TodoService) CompleteTask(id string) error {

// }
