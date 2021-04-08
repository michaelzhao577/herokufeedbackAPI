package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Response struct {
	gorm.Model

	Service  string
	Rating   int
	Feedback string
}

var db *gorm.DB
var err error

func main() {
	// open connection to db
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database")
	}

	// close connection to db when main func terminates
	defer db.Close()

	// make migration to the db if they have not already been created
	db.AutoMigrate(&Response{})

	router := handleRequests()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func getResponses(w http.ResponseWriter, r *http.Request) {
	var responses []Response
	db.Find(&responses)
	json.NewEncoder(w).Encode(responses)
}

func createResponse(w http.ResponseWriter, r *http.Request) {
	var response Response
	json.NewDecoder(r.Body).Decode(&response)

	createdResponse := db.Create(&response)
	err = createdResponse.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&response)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "homepage endpoint reached")
}

func handleRequests() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", homePage)
	router.HandleFunc("/responses", getResponses).Methods("GET")
	router.HandleFunc("/create/response", createResponse).Methods("POST")

	return router
}
