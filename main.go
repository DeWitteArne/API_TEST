package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday"})
	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person

func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
		}
	}
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	var id = strconv.Itoa(len(people) + 1)
	person.ID = string(id)
	people = append(people, person)
	json.NewEncoder(w).Encode("Success")
}
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var deleted = false
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			deleted = true
			break
		}
	}

	if deleted {
		json.NewEncoder(w).Encode("Success")
	} else {
		json.NewEncoder(w).Encode("There was nothing to delete")

	}
}
