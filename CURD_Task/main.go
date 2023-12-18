package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rajashri570/Go_PROJECT/CURD_Task/Task"
)

func initializeRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/task", Task.View_tasks).Methods("GET")
	r.HandleFunc("/task/{id}", Task.Get_task).Methods("GET")
	r.HandleFunc("/task", Task.Create_task).Methods("POST")
	// r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	// r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
	fmt.Println("server started..")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	Task.InitialMigration()
	initializeRouter()
}
