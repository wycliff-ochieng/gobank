package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wycliff-ochieng/data"
)

type APIServer struct {
	Addr  string
	Store data.Storage
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

func NewAPIServer(addr string, store data.Storage) *APIServer {
	//	store, err := data.NewPostgrestore()
	//	if err != nil {
	//		panic(err)
	//	}
	//	if err := store.Init(); err != nil {
	//		panic(err)
	//	}
	return &APIServer{
		Addr:  addr,
		Store: store,
	}
}

// handlefunc takes a path, http handler
func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", MakeHttpHandlerFunc(s.handleAccount))
	//router.HandleFunc("/",MakeHttpHandlerFunc(s.handleGetAccount))
	router.HandleFunc("/account/{id}", MakeHttpHandlerFunc(s.handleGetAccountByID))
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
	accounts, err := s.Store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)["id"]

	id, err := strconv.Atoi(vars)
	if err != nil {
		return err
	}

	account, err := s.Store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	createaccountrequest := new(data.CreateAccountReq)

	if err := json.NewDecoder(r.Body).Decode(createaccountrequest); err != nil {
		return err
	}

	account := data.NewAccount(createaccountrequest.Firstname, createaccountrequest.Lastname)

	if err := s.Store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransferAccount() error {
	return nil
}
