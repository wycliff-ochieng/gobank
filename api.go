package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

// serializing
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type APIFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

// make api func into an http handler func
func MakeHttpHandler(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

type APIserver struct {
	ListenAddr string
}

func NewAPIServer(ListenAddr string) *APIserver {
	return &APIserver{
		ListenAddr: ListenAddr,
	}
}

func (s *APIserver) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", MakeHttpHandler(s.handleAccount))
	http.ListenAndServe(s.ListenAddr, router)
}

func (s *APIserver) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "Delete" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method does not exist %s", r.Method)
}

func (s *APIserver) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("John", "Mapangale")
	return WriteJSON(w, http.StatusOK, account)
}
func (s *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Account struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Number    int    `json:"number"`
	Balance   int64  `json:"balance`
}

func NewAccount(firstname, lastname string) *Account {
	return &Account{
		ID:        rand.Intn(10),
		Firstname: firstname,
		Lastname:  lastname,
		Number:    rand.Intn(1000000),
	}
}

func main() {
	fmt.Println("Hawayu")
	server := NewAPIServer(":3000")
	server.Run()
}
