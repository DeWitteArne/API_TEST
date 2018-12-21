package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var people = []*Person{
	NewPerson("John", "Doe", NewAddress("City x", "State x")),
	NewPerson("Mary", "The Lamb", NewAddress("City a", "State x")),
	NewPerson("Kees", "De Jong", NewAddress("City a", "State 3")),
	NewPerson("Oester", "Fish", NewAddress("City x", "State d")),
	NewPerson("Poison", "Dart", NewAddress("City x", "State i")),
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/people", makeHTTPHandler(GetPeople)).Methods("GET")
	r.HandleFunc("/people/{id}", makeHTTPHandler(GetPerson)).Methods("GET")
	r.HandleFunc("/people", makeHTTPHandler(CreatePerson)).Methods("POST")
	r.HandleFunc("/people/{id}", makeHTTPHandler(DeletePerson)).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}

type httpAPIFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandler(f httpAPIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			w.Write([]byte(err.Error()))
		}
	}
}

// Person represents a person.
type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

// NewPerson returns a new person object.
func NewPerson(firstname string, lastname string, address *Address) *Person {
	return &Person{
		Firstname: firstname,
		Lastname:  lastname,
		Address:   address,
	}
}

// Address holds address details
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

// NewAddress returns a new address object.
func NewAddress(city, state string) *Address {
	return &Address{
		City:  city,
		State: state,
	}
}

// GetPeople returns a list of persons.
func GetPeople(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, people)
}

// GetPerson returns a person by it's ID.
func GetPerson(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			return writeJSON(w, http.StatusOK, item)
		}
	}
	return nil
}

// CreatePerson creates a new person.
func CreatePerson(w http.ResponseWriter, r *http.Request) error {
	person := Person{ID: strconv.Itoa(len(people) + 1)}

	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		return err
	}
	defer r.Body.Close()

	people = append(people, &person)

	return writeJSON(w, http.StatusOK, "success")
}

// DeletePerson deletes a person by it's ID.
func DeletePerson(w http.ResponseWriter, r *http.Request) error {
	var (
		params  = mux.Vars(r)
		deleted = false
	)

	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			deleted = true
			break
		}
	}

	var msg string
	if deleted {
		msg = "Success"
	} else {
		msg = "There was nothing to delete"
	}

	return writeJSON(w, http.StatusOK, msg)
}

func writeJSON(w http.ResponseWriter, i int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(i)

	return json.NewEncoder(w).Encode(v)
}
