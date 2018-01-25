package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("Starting client...")

	//client example
	resp, err := http.Get("http://localhost:8888/endpoint1")
	if err != nil {
		//handle error
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
}
