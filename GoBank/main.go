package main

func main() {
	srv := NewAPIServer(":8080")
	srv.Run()
}
