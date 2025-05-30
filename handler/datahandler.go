package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
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
	router.HandleFunc("/account/{id}", withJWTAuth(MakeHttpHandlerFunc(s.handleGetAccountByID)))
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
	if r.Method == "GET" {
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

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method no allowed %s", r.Method)
	//	vars := mux.Vars(r)["id"]
	//
	//	id, err := strconv.Atoi(vars)
	//	if err != nil {
	//		return err
	//	}
	//
	//	account, err := s.Store.GetAccountByID(id)
	//	if err != nil {
	//		return err
	//	}

	// return WriteJSON(w, http.StatusOK, account)
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

	tokenString, err := CreateJWT(account)
	if err != nil {
		return err
	}

	fmt.Println("JWT token", tokenString)

	return WriteJSON(w, http.StatusCreated, account)

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.Store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransferAccount() error {
	return nil
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given")
	}
	return id, nil
}

//JWT handling and fixing

func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling jwt function")

		tokenString := r.Header.Get("x-jwt-token")
		_, err := validateJWT(tokenString)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, APIError{Error: "invalid token"})
			return
		}
		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := "mystupidsecret"

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func CreateJWT(account *data.Account) (string, error) {
	claims := jwt.MapClaims{
		"expiresAt":     150000,
		"accountNumber": account.ID,
	}
	secret := "mystupidsecret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
