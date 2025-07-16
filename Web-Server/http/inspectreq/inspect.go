package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.TODO()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://example.com/index.html", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "deflate")
	req.Header.Set("User-Agent", "eblog/1.0")
	req.Header.Set("some-key", "a value")
	req.Header.Set("SOMe-KEY", "somevalue") // overwrites previous

	req.Write(os.Stdout) // print full HTTP request
}
