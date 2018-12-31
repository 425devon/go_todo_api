package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/425devon/go_todo_api/pkg/mongo"
	"github.com/425devon/go_todo_api/pkg/server"
)

func main() {
	ms, err := mongo.NewSession("127.0.0.1:27017")
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}
	u := mongo.NewTodoService(ms.Copy(), "go_todo_server", "todo")
	s := server.NewServer(u)

	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, os.Interrupt)

	go func() {
		sig := <-gracefulStop
		fmt.Printf("\ncaught sig: %+v\n", sig)
		//dropping db during testing
		fmt.Println("Dropping and closing db")
		ms.DropDatabase("go_todo_server")
		//
		ms.Close()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	s.Start()
}
