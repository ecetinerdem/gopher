package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type api struct {
	addr string
}

var users = []User{}

func (a *api) createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = insertUser(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *api) getUsersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	err := encoder.Encode(users)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/":
			w.Write([]byte("index page"))
			return
		case "/users":
			w.Write([]byte("users page"))
			return
		}
	default:
		w.Write([]byte("users page"))
		return

	}

}

func insertUser(u User) error {
	if u.FirstName == "" {
		return errors.New("First name cannot be missing")
	}

	if u.LastName == "" {
		return errors.New("Last name cannot be missing")
	}

	for _, user := range users {
		if user.FirstName == u.FirstName && user.LastName == u.LastName {
			return errors.New("User already exist")
		}
	}
	users = append(users, u)
	return nil
}
