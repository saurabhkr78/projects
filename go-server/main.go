package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	} else if r.Method != http.MethodGet {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	} else {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "world"
		}
		fmt.Fprintf(w, "hello %s", name)
	}
}
func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/form" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	} else if r.Method != http.MethodPost {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form data", http.StatusBadRequest)
			return
		}
		name := r.FormValue("name")
		email := r.FormValue("email")

		fmt.Fprintf(w, "Hello %s, your email is %s", name, email)
	}
}
func main() {
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", fs)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("server is started on 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
