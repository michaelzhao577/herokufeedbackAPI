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
	// db connection string
	// dbURI := "postgres://mrefuzdezkdqil:d07ad6986d8d15985a35cf47ab018d1d6bd228f9c42fa75f97f12b8a606e5b8a@ec2-3-233-43-103.compute-1.amazonaws.com:5432/d2vcvaemk0fbjn"

	// open connection to db
	dbURL := os.Getenv("DATABASE_URL")
	dbURL = dbURL + "sslmode=require"
	db, err = gorm.Open("postgres", dbURL)
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

	log.Fatal(http.ListenAndServe(":80", router))
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

func handleRequests() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/responses", getResponses).Methods("GET")
	router.HandleFunc("/create/response", createResponse).Methods("POST")

	return router
}
