package server

import (
	"log"
	"net/http"
	"os"

	"github.com/425devon/go_todo_api/pkg/mongo"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer(s *mongo.TodoService) *Server {
	server := Server{router: mux.NewRouter()}
	NewTodoRouter(s, server.newSubRouter("/todo"))
	return &server
}

func (server *Server) Start() {
	log.Println("Server starting on port :8080")
	if err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, server.router)); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}

func (server *Server) newSubRouter(path string) *mux.Router {
	return server.router.PathPrefix(path).Subrouter()
}
