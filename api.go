package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddress string
	storage       Storage
}

func NewApiServer(listenAddress string, storage Storage) *ApiServer {
	return &ApiServer{
		listenAddress: listenAddress,
		storage:       storage,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", s.withJWTAuth(makeHTTPHandleFunc(s.handleAccountById)))

	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("JSON API server running on port")

	http.ListenAndServe(s.listenAddress, router)
}

func (s *ApiServer) handleLogin(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return fmt.Errorf("method not allowed: %s", r.Method)
	}

	loginReq := new(LoginRequest)
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(loginReq); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, loginReq)
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed: %s", r.Method)
}

func (s *ApiServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccountById(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed: %s", r.Method)
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.storage.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)

	if err != nil {
		return err
	}

	account, err := s.storage.GetAccountById(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	createAccountReq := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account, err := NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Password)

	if err != nil {
		return nil
	}

	newAccount, err := s.storage.CreateAcccount(account)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, newAccount)
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.storage.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)

	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}

	return id, nil
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{"Invalid Token"})
}

func (s *ApiServer) withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)
		if err != nil || !token.Valid {
			permissionDenied(w)
			return
		}

		accountId, err := getID(r)

		if err != nil {
			permissionDenied(w)
			return
		}

		account, err := s.storage.GetAccountById(accountId)

		if err != nil {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

const secret = "someSecret"

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjoyMzE1OTAsImV4cGlyZXNBdCI6MTUwMH0.dUngXP_QSeaHuaDs5fUtUTpedXvpSNeg3VBO2LKZdA8
func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     1500,
		"accountNumber": account.Number,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
