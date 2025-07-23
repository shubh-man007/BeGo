package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	storage    *PostgresStorage
}

func NewAPIServer(listenAddr string, storage *PostgresStorage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))
	router.HandleFunc("/account/{id}/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("Server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)

}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)

	case "POST":
		return s.handleCreateAccount(w, r)

	case "DELETE":
		return s.handleDeleteAccount(w, r)

	case "PUT":
		return s.handleUpdateAccount(w, r)

	default:
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid account id")
	}

	p := s.storage
	acc, err := p.GetAccountByID(id)
	if err != nil {
		log.Printf("Could not fetch account[%d]", id)
		return err
	}
	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Could not read request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError) // HTTP 500
		return err
	}
	defer r.Body.Close()

	var a account
	err = json.Unmarshal(body, &a)
	if err != nil {
		log.Printf("Could not unmarshal request body: %v", err)
		w.WriteHeader(http.StatusBadRequest) // HTTP 400
		return err
	}

	p := s.storage
	err = p.CreateAccount(&a)
	if err != nil {
		log.Printf("Could not create account[%d]", a.AccountID)
		return err
	}
	return WriteJSON(w, http.StatusOK, a)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid account id")
	}

	p := s.storage
	err = p.DeleteAccount(id)
	if err != nil {
		log.Printf("Could not delete account[%d]", id)
		return err
	}

	del := "Account" + idStr + " has been deleted"
	return WriteJSON(w, http.StatusOK, del)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Could not read request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError) // HTTP 500
		return err
	}

	defer r.Body.Close()

	var a account
	err = json.Unmarshal(body, &a)
	if err != nil {
		log.Printf("Could not unmarshal request body: %v", err)
		w.WriteHeader(http.StatusBadRequest) // HTTP 400
		return err
	}

	p := s.storage
	err = p.UpdateAccount(&a)
	if err != nil {
		log.Printf("Could not update account[%d]", a.AccountID)
		return err
	}

	return WriteJSON(w, http.StatusOK, a)

}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid account id")
	}

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("invalid request body: %v", err)
	}
	defer r.Body.Close()

	acc, err := s.storage.GetAccountByID(id)
	if err != nil {
		return fmt.Errorf("account not found: %v", err)
	}

	updatedAcc, err := s.storage.TransferBalance(acc, req.Amount)
	if err != nil {
		return fmt.Errorf("could not update balance: %v", err)
	}

	return WriteJSON(w, http.StatusOK, updatedAcc)
}
