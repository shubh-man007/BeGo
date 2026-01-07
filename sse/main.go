package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const eventsHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Stream Events</title>
</head>
<body>
    <h2>Model Response</h2>
    <p id="resp"></p>

    <script>
        var sse = new EventSource('/events');
        var resp = document.getElementById('resp')

        sse.onmessage = function (event) {
            resp.innerHTML += event.data + ' ';
        };

		sse.onerror = function (e) {
            setTimeout(() => (resp.innerHTML = ''), 1000)
        };
    </script>
</body>
</html>
`

const message = `
<html>
<head>
	<title>About</title>
</head>
<body>
	<p><stronig>Hi, <i>waddup</i> !</strong> This is a page about <b>SSE</b><p>
</body>
</html>
`

type ModelRes struct {
	output []string
}

func (tk *ModelRes) streamResHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.WriteHeader(200)

	for _, val := range tk.output {
		content := fmt.Sprintf("data: %s\n\n", string(val))
		w.Write([]byte(content))
		w.(http.Flusher).Flush()

		time.Sleep(time.Millisecond * 500)
	}
}

func main() {
	mux := http.NewServeMux()

	tk := ModelRes{
		output: []string{"Am", "I", "a", "bloody", "AI", "generated", "response", "?"},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		w.Write([]byte(eventsHTML))
	})

	mux.HandleFunc("/events", tk.streamResHandler)

	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)

		_, err := w.Write([]byte(message))
		if err != nil {
			log.Print("Could not write to about")
		}
	})

	s := &http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	log.Print("Listening at port :5000")
	err := s.ListenAndServe()
	if err != nil {
		log.Printf("Could not listen to server at port: 5000. Err: %v", err)
	}
}
