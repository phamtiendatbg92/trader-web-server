package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// Book struct model
type Book struct {
	Id     string  `json:"id"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}
type Author struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

func getAllTutorial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tutorial := client.Database("trader-web-mongodb").Collection("tutorials")
	log.Print("11111")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, currErr := tutorial.Find(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	log.Print("22222")
	defer cur.Close(ctx)

	var posts []Tutorial
	err := cur.All(ctx, &posts)
	if err != nil {
		log.Print("333")
		panic(err)
	}

	json.NewEncoder(w).Encode(posts)
}

func getall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.Id = strconv.Itoa(currentID)
	currentID++
	books = append(books, book)

	json.NewEncoder(w).Encode(books)
}
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get params
	for _, value := range books {
		if value.Id == params["id"] {
			json.NewEncoder(w).Encode(value)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(books)
}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

var books []Book
var currentID int = 3

func main() {
	connect()
	books = append(books, Book{Id: "1", Title: "title 1", Author: &Author{Name: "AAA", Age: 23}})
	books = append(books, Book{Id: "2", Title: "title 2", Author: &Author{Name: "BBB", Age: 23}})
	// Create router
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/get-all-tutorial", getAllTutorial).Methods("GET")
	router.HandleFunc("/api/v1/getall", getall).Methods("GET")
	router.HandleFunc("/api/v1/get-book/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/v1/add-book", addBook).Methods("POST")
	router.HandleFunc("/api/v1/update-book", updateBook).Methods("PUT")
	router.HandleFunc("/api/v1/delete-book", deleteBook).Methods("DELETE")

	// start Server
	log.Fatal(http.ListenAndServe(":5000", router))
}
