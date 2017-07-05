package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"./api"
	"./model"
)

func main() {
	log.Println("Connecting to rethinkDB on localhost...")
	err := model.InitSession()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	r := mux.NewRouter()
	r.HandleFunc("/search/{query}", api.SearchDocumentsHandler).Methods("GET")

	log.Println("Running the server on 8000...")
	http.ListenAndServe(":8000", r)
}