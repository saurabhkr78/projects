package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type movie struct {
	ID       *string   `json:"id"`
	Title    *string   `json:"title"`
	Year     *string   `json:"year"`
	Director *director `json:"director"`
}
type director struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

// data base slice of movies
var Movie []movie

//create all the handlers of the routes

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(Movie)
}
func getMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	for _, movie := range Movie {
		if movie.ID != nil && *movie.ID == id {

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(movie); err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}

			return
		}
	}

	http.Error(w, "movie not found ", http.StatusNotFound)
}
func createMovie(w http.ResponseWriter, r *http.Request) {
	//create a new movie variable where the data just read from the request body will be stored
	var newMovie movie
	if err := json.NewDecoder(r.Body).Decode(&newMovie); err != nil {
		http.Error(w, "failed to decode req body", http.StatusBadRequest)
		return
	}
	// duplicate check
	for _, m := range Movie {
		if m.Title == newMovie.Title && m.Year == newMovie.Year {
			http.Error(w, "movie already exists", http.StatusConflict)
			return
		}
	}
	//generate a random id for the movie
	newMovie.ID = new(string)
	*newMovie.ID = strconv.Itoa(rand.Intn(1000000))
	//append the new movie to the database
	Movie = append(Movie, newMovie)
	//return the new movie created in the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newMovie); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}

}
func updateMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	var updatedMovie movie

	if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	for i, m := range Movie {
		if m.ID != nil && *m.ID == id {

			// top-level fields
			if updatedMovie.Title != nil {
				Movie[i].Title = updatedMovie.Title
			}

			if updatedMovie.Year != nil {
				Movie[i].Year = updatedMovie.Year
			}

			// nested struct
			if updatedMovie.Director != nil {

				if Movie[i].Director == nil {
					Movie[i].Director = &director{}
				}

				if updatedMovie.Director.FirstName != nil {
					Movie[i].Director.FirstName = updatedMovie.Director.FirstName
				}

				if updatedMovie.Director.LastName != nil {
					Movie[i].Director.LastName = updatedMovie.Director.LastName
				}
			}

			json.NewEncoder(w).Encode(Movie[i])
			return
		}
	}

	http.Error(w, "movie not found", http.StatusNotFound)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	for i, m := range Movie {
		if m.ID != nil && *m.ID == id { //since m.ID is * string so can be nil and also we need to dereference it to get the value to compare with id
			Movie = append(Movie[:i], Movie[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{"message": "movie deleted successfully"})
			return
		}
	}

	http.Error(w, "movie not found", http.StatusNotFound)
}

func main() {
	//craete a hardcoded movie database
	Movie = append(Movie, movie{ID: new(string), Title: new(string), Year: new(string), Director: &director{FirstName: new(string), LastName: new(string)}})
	*Movie[0].ID = "1"
	*Movie[0].Title = "The Matrix"
	*Movie[0].Year = "1999"
	*Movie[0].Director.FirstName = "Lana"
	*Movie[0].Director.LastName = "Wachowski"

	Movie = append(Movie, movie{ID: new(string), Title: new(string), Year: new(string), Director: &director{FirstName: new(string), LastName: new(string)}})
	*Movie[1].ID = "2"
	*Movie[1].Title = "Inception"
	*Movie[1].Year = "2010"
	*Movie[1].Director.FirstName = "Christopher"
	*Movie[1].Director.LastName = "Nolan"

	Movie = append(Movie, movie{ID: new(string), Title: new(string), Year: new(string), Director: &director{FirstName: new(string), LastName: new(string)}})
	*Movie[2].ID = "3"
	*Movie[2].Title = "Interstellar"
	*Movie[2].Year = "2014"
	*Movie[2].Director.FirstName = "Christopher"
	*Movie[2].Director.LastName = "Nolan"

	Movie = append(Movie, movie{ID: new(string), Title: new(string), Year: new(string), Director: &director{FirstName: new(string), LastName: new(string)}})
	*Movie[3].ID = "4"
	*Movie[3].Title = "The Dark Knight"
	*Movie[3].Year = "2008"
	*Movie[3].Director.FirstName = "Christopher"
	*Movie[3].Director.LastName = "Nolan"

	Movie = append(Movie, movie{ID: new(string), Title: new(string), Year: new(string), Director: &director{FirstName: new(string), LastName: new(string)}})
	*Movie[4].ID = "5"
	*Movie[4].Title = "Pulp Fiction"
	*Movie[4].Year = "1994"
	*Movie[4].Director.FirstName = "Quentin"
	*Movie[4].Director.LastName = "Tarantino"

	//create router
	r := mux.NewRouter()

	//regsister the handlers for the routes
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
