package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/425devon/go_todo_api/pkg/models"
	"github.com/425devon/go_todo_api/pkg/mongo"
	"github.com/gorilla/mux"
)

type todoRouter struct {
	todoService *mongo.TodoService
}

func NewTodoRouter(s *mongo.TodoService, router *mux.Router) *mux.Router {
	todoRouter := todoRouter{s}
	router.HandleFunc("/", todoRouter.welcomeTest).Methods("GET")
	router.HandleFunc("/lists", todoRouter.getAllListsHandler).Methods("GET")
	router.HandleFunc("/lists", todoRouter.createListHandler).Methods("POST")
	router.HandleFunc("/list/{id}/tasks", todoRouter.createTaskHandler).Methods("POST")
	return router
}

func (tr *todoRouter) welcomeTest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Welcome to go-todo")
}

func (tr *todoRouter) createListHandler(w http.ResponseWriter, r *http.Request) {
	list, err := decodeList(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	_, err = tr.todoService.CreateList(list)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Json(w, http.StatusCreated, err)
}

func (tr *todoRouter) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := decodeTask(r)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	list, err := tr.todoService.GetListByID(id)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = tr.todoService.CreateTask(&list, task)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Json(w, http.StatusCreated, err)
}

func (tr *todoRouter) getAllListsHandler(w http.ResponseWriter, r *http.Request) {
	lists, err := tr.todoService.GetAllLists()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Json(w, http.StatusOK, lists)
}

func (tr *todoRouter) getListByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	list, err := tr.todoService.GetListByID(id)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
	}
	Json(w, http.StatusOK, list)
}
func decodeList(r *http.Request) (*models.TodoList, error) {
	var list models.TodoList
	if r.Body == nil {
		return nil, errors.New("no request body")
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&list)
	return &list, err
}

func decodeTask(r *http.Request) (*models.Task, error) {
	var task models.Task
	if r.Body == nil {
		return nil, errors.New("no request body")
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	return &task, err
}
