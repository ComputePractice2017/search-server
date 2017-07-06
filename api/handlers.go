package api

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/ComputePractice2017/search-server/model"
	"github.com/gorilla/mux"
)

func SearchDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	
	persons, err := model.FindDocs(vars["query"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := json.NewEncoder(w).Encode(persons); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}