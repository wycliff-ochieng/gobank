package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wycliff-ochieng/data"
)

type APIServer struct {
	Addr    string
	Storage data.Storage
}

// function signature of function we are using
type apiFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Error string
}

// decorate our apifunc into http handler func
func MakeHttpHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusConflict, APIError{Error: err.Error()})
		}
	}
}

// serialize the input
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func NewAPIServer(addr string) *APIServer {
	store, err := data.NewPostgrestore()
	if err != nil {
		panic(err)
	}
	if err := store.Init(); err != nil {
		panic(err)
	}
	return &APIServer{
		Addr:    addr,
		Storage: store,
	}
}

// handlefunc takes a path, http handler
func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/", MakeHttpHandlerFunc(s.handleAccount))
	http.ListenAndServe(s.Addr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not allowed")
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account := data.NewAccount("Jessy", "Elikanah")
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Firstname string
		Lastname  string
		Balance   int
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("Couldnt unmarshal: %v", err)
	}

	account := data.NewAccount(req.Firstname, req.Lastname)
	account.Balance = req.Balance
	if err := s.Storage.CreateAccount(account); err != nil {
		return fmt.Errorf("Couldnt create account : %v", err)
	}

	return WriteJSON(w, http.StatusCreated, account)

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransferAccount() error {
	return nil
}
