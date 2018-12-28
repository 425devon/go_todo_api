package mongo_test

import (
	"log"
	"testing"

	"github.com/425devon/go_todo_api/pkg/models"
	"github.com/425devon/go_todo_api/pkg/mongo"
)

const (
	mongoURL           = "localhost:27017"
	dbName             = "todo_test_db"
	todoCollectionName = "todo"
)

func Test_todoService(t *testing.T) {
	t.Run("CreateList", createList_should_create_list_and_return_id)
	t.Run("GetListByID", getListByID_should_get_list_by_ID)
	t.Run("CreateTask", createTask_should_create_task_and_add_to_list)
	t.Run("GetTaskByID", getTaskByID_should_find_task_by_id)
}

func createList_should_create_list_and_return_id(t *testing.T) {
	//Arrange
	session := newSession()
	todoService := newTodoService(session)
	defer dropAndCloseDB(session)

	todoList := models.TodoList{
		Name:        "test_list",
		Description: "this list is for testing",
		Tasks:       nil,
	}

	//Act
	uid, err := todoService.CreateList(&todoList)

	//Assert
	if err != nil {
		t.Errorf("Unable to create list: `%s`", err)
	}
	if len(uid) == 0 {
		t.Errorf("Expected list ID, Got: `%s`", uid)
	}
}

func getListByID_should_get_list_by_ID(t *testing.T) {
	//Arrange
	session := newSession()
	todoService := newTodoService(session)
	defer dropAndCloseDB(session)

	todoList := models.TodoList{
		Name:        "test_list",
		Description: "this list is for testing",
		Tasks:       nil,
	}

	//Act
	uid, _ := todoService.CreateList(&todoList)
	list, err := todoService.GetListByID(uid)

	//Assert
	if err != nil {
		t.Errorf("Unable to get list: `%s`", err)
	}
	if list.Name != "test_list" {
		t.Errorf("Incorrect list name expected: `test_list` got: `%s`", list.Name)
	}
}

func createTask_should_create_task_and_add_to_list(t *testing.T) {
	//Arrange
	session := newSession()
	todoService := newTodoService(session)
	defer dropAndCloseDB(session)

	todoList := models.TodoList{
		Name:        "test_list",
		Description: "this list is for testing",
		Tasks:       nil,
	}

	task := models.Task{
		Name: "test_task",
	}

	//Act
	uid, _ := todoService.CreateList(&todoList)
	list, _ := todoService.GetListByID(uid)
	tid, err := todoService.CreateTask(&list, &task)

	//Assert
	if err != nil {
		t.Errorf("Unable to create task: `%s`", err)
	}
	if len(tid) == 0 {
		t.Errorf("Expected Task ID, Got: `%s`", tid)
	}

}

func getTaskByID_should_find_task_by_id(t *testing.T) {
	//Arrange
	session := newSession()
	todoService := newTodoService(session)
	defer dropAndCloseDB(session)

	todoList := models.TodoList{
		Name:        "test_list",
		Description: "this list is for testing",
		Tasks:       nil,
	}

	task := models.Task{
		Name: "test_task",
	}

	//Act
	uid, _ := todoService.CreateList(&todoList)
	list, _ := todoService.GetListByID(uid)
	tid, _ := todoService.CreateTask(&list, &task)
	tsk, err := todoService.GetTaskByID(tid)

	//Assert
	if err != nil {
		t.Errorf("Unable to retrieve task: `%s`", err)
	}
	if tsk.Completed != false {
		t.Errorf("Expected Completed status to be `false` got: `%v`", tsk.Completed)
	}
}

func newSession() *mongo.Session {
	session, err := mongo.NewSession(mongoURL)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	return session
}

func newTodoService(session *mongo.Session) *mongo.TodoService {
	todoService := mongo.NewTodoService(session.Copy(), dbName, todoCollectionName)
	return todoService
}

func dropAndCloseDB(session *mongo.Session) {
	session.DropDatabase(dbName)
	session.Close()
}