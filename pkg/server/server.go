package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var servehtml = `
<!DOCTYPE html>
<html>
	<body>
		<h1 style="text-align:center;" > User Database </h1>
		<p style="text-align:center;" > Welcome to the user database. </p>
	</body>
</html>
`

// user represents the JSON value that's sent as a response to a request.
type user struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// userinfo is the information that is stored per user.
type userinfo struct {
	email string
	age   int
}

type server struct {
	users map[string]userinfo
}

func New() *server {
	return &server{
		users: make(map[string]userinfo),
	}
}

func (s *server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(servehtml))
}

func (s *server) HandleCreateUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Could not read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError) // HTTP 500
			return
		}
		defer r.Body.Close()

		var u user
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Printf("Could not unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest) // HTTP 400
			return
		}

		log.Printf("Create User: %v", u.Name)
		s.users[u.Name] = userinfo{
			email: u.Email,
			age:   u.Age,
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed) // HTTP 415
	}
}

func (s *server) HandleUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	// name := r.URL.Query().Get("name")
	res, ok := s.users[name]
	if !ok {
		log.Printf("User: %s does not exist\n", name)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		ret := user{
			Name:  name,
			Email: res.email,
			Age:   res.age,
		}
		msg, err := json.Marshal(ret)
		if err != nil {
			log.Printf("Unalble to marshal server response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("Get user %s", name)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(msg)

	case http.MethodPatch:
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Could not read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError) // HTTP 500
			return
		}

		defer r.Body.Close()

		log.Printf("Update user %s", name)
		var u user
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Printf("Could not unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest) // HTTP 400
			return
		}

		updated := s.users[name]
		if u.Age != 0 {
			updated.age = u.Age
		}
		if u.Email != "" {
			updated.email = u.Email
		}

		s.users[name] = updated

	case http.MethodDelete:
		log.Printf("Delete user %s", name)
		delete(s.users, name)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
