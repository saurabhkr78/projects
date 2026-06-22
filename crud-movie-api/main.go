package crudmovieapi

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type movie struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Year     string    `json:"year"`
	Director *director `json:"director"`
}
type director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// data base slice of movies
var Movie []movie

//create all the handlers of the routes

func getMovies(w http.ResponseWriter, r *http.Request) {

}
func getMovie(w http.ResponseWriter, r *http.Request)    {}
func createMovie(w http.ResponseWriter, r *http.Request) {}
func updateMovie(w http.ResponseWriter, r *http.Request) {}
func deleteMovie(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := mux.NewRouter()
	//craete a hardcoded movie database
	Movie = append(Movie, movie{ID: "1", Title: "The Matrix", Year: "1999", Director: &director{FirstName: "Lana", LastName: "Wachowski"}})
	Movie = append(Movie, movie{ID: "2", Title: "Inception", Year: "2010", Director: &director{FirstName: "Christopher", LastName: "Nolan"}})
	Movie = append(Movie, movie{ID: "3", Title: "Interstellar", Year: "2014", Director: &director{FirstName: "Christopher", LastName: "Nolan"}})
	Movie = append(Movie, movie{ID: "4", Title: "The Dark Knight", Year: "2008", Director: &director{FirstName: "Christopher", LastName: "Nolan"}})
	Movie = append(Movie, movie{ID: "5", Title: "Pulp Fiction", Year: "1994", Director: &director{FirstName: "Quentin", LastName: "Tarantino"}})
	//get all
	r.HandleFunc("/movies", getMovies).Methods("GET")

	//get by id->getMovie
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	// create ->createMovie
	r.HandleFunc("/movies", createMovie).Methods("POST")
	//update ->updateMovie
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	// //delete->deleteMovie
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	//start the server
	fmt.Println("starting the server @8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
