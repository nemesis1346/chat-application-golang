package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")
	h := http.NewServeMux()

	h.HandleFunc("/endpoint1", getMessage)

	err := http.ListenAndServe(":8888", h)
	log.Fatal(err)
}

//My first endpoint
func getMessage(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world")
}
