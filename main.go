package main

import (
	"log"
	"net/http"
	"service/app/controller"
	"service/database"

	"github.com/gorilla/mux"
)

func init() {
	_, err := database.DatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}

}
func main() {

	router := mux.NewRouter()
	router.HandleFunc("/create", controller.CreateBook).Methods(http.MethodPost)
	router.HandleFunc("/book/{id}", controller.GetBookById).Methods(http.MethodGet)
	router.HandleFunc("/book/{id}", controller.UpdateBookById).Methods(http.MethodPut)
	router.HandleFunc("/book/{id}", controller.DeleteBookById).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", router))
}
