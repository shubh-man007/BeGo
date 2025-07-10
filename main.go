// package main

// import (
// 	"log"
// 	"net/http"
// 	"time"
// )

// var servehtml = `
// <!DOCTYPE html>
// <html>
// <head>
// 	<title>Welcome</title>
// 	<style>
// 		body {
// 			background-color:rgb(250, 192, 192);
// 			font-family: Arial, sans-serif;
// 			text-align: center;
// 			padding-top: 50px;
// 		}
// 		h1 {
// 			color: #333;
// 		}
// 		p {
// 			color: #666;
// 		}
// 	</style>
// </head>
// <body>
// 	<h1>Welcome to My Go Web Server!</h1>
// 	<p>This page is served using Go's <code>net/http</code> package.</p>
// </body>
// </html>
// `

// var UserInfo = `{
// 	"UserId" : 123,
// 	"UserName" : "Shubh"
// }
// `

// func main() {
// 	address := ":8080"
// 	// Multiplexer is used for routing
// 	// A ServeMux matches incoming requests to registered handler functions based on URL paths.
// 	mux := http.NewServeMux()

// 	// w http.ResponseWriter: Used to send the response back to the client.
// 	// r *http.Request: Contains all the info about the incoming request.
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		// for key, values := range r.Header {
// 		// 	for _, value := range values {
// 		// 		w.Write([]byte(key + ":" + value + "\n"))
// 		// 	}
// 		// }

// 		w.Header().Add("content-type", "text/html")

// 		// for key, values := range w.Header() {
// 		// 	for _, value := range values {
// 		// 		w.Write([]byte(key + ":" + value + "\n"))
// 		// 	}
// 		// }
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(servehtml))
// 	})

// 	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("content-type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(UserInfo))
// 	})

// 	s := &http.Server{
// 		Addr:           address,
// 		Handler:        mux,
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 	}

// 	log.Printf("Server running at %v", address)
// 	log.Fatal(s.ListenAndServe())
// }

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/shubh-man007/BeGo/pkg/server"
)

func main() {
	address := ":8080"
	srv := server.New()

	mux := http.NewServeMux()

	mux.HandleFunc("/", srv.HandleIndex)
	mux.HandleFunc("/userCreate", srv.HandleCreateUsers)

	s := &http.Server{
		Addr:           address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Server running at %v", address)
	log.Fatal(s.ListenAndServe())
}
