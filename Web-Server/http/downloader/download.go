package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	dir     = flag.String("dir", ".", "Mention directory to be used.")
	timeout = flag.Duration("timeout", 30*time.Second, "Timeout for the HTTP request.")
)

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatal("usage: download [-timeout duration] [-dir path] url filename")
	}

	url, filename := flag.Arg(0), flag.Arg(1)

	client := &http.Client{Timeout: *timeout}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	if err := downloadAndSave(ctx, client, url, *dir, filename); err != nil {
		log.Fatalf("download failed: %v", err)
	}

	fmt.Printf("Downloaded file saved to: %s\n", filepath.Join(*dir, filename))
}

func downloadAndSave(ctx context.Context, c *http.Client, url, dir, filename string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: GET %q: %v", url, err)
	}

	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("response status: %s", res.Status)
	}

	dstPath := filepath.Join(dir, filename)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("creating file: %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, res.Body); err != nil {
		return fmt.Errorf("copying response to file: %v", err)
	}

	return nil
}

// go run download.go <https://site-name/> output.html
