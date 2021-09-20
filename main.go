package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"trader-web-api/dbcontroller"

	"github.com/gorilla/mux"
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
	//db - controller.connect()
	dbcontroller.Connect()
	books = append(books, Book{Id: "1", Title: "title 1", Author: &Author{Name: "AAA", Age: 23}})
	books = append(books, Book{Id: "2", Title: "title 2", Author: &Author{Name: "BBB", Age: 23}})
	// Create router
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/getall", getall).Methods("GET")
	router.HandleFunc("/api/v1/get-book/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/v1/add-book", addBook).Methods("POST")
	router.HandleFunc("/api/v1/update-book", updateBook).Methods("PUT")
	router.HandleFunc("/api/v1/delete-book", deleteBook).Methods("DELETE")

	router.HandleFunc("/api/v1/get-list-tutorials", getAllTutorial).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/detail-tutorial/{url}", getDetailTutorial).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/get-hashtag", getAllHashTag).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/upload-new-post", uploadNewPost).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/update-post", updateTutorial).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/v1/delete-post/{id}", deleteTutorial).Methods(http.MethodDelete, http.MethodOptions)

	// router.PathPrefix("/images").Handler(http.FileServer(http.Dir("./public/images")))

	// This will serve files under http://localhost:8000/static/<filename>
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./public/images/"))))

	router.Use(accessControlMiddleware)
	// start Server
	log.Fatal(http.ListenAndServe(":5000", router))
}

// access control and  CORS middleware
func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
