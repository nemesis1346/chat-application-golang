package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	. "testing"
)

//type message struct{}
type test int

func main() {
	fmt.Println("Starting server...")
	h := http.NewServeMux()

	//DEFINITION OF METHODS
	h.HandleFunc("/endpoint1", createChatRoom)
	h.HandleFunc("/createChatRoom", createChatRoom)
	h.HandleFunc("/listChatRoom", listChatRoom)
	h.HandleFunc("/joinChatRoom", jointChatRoom)
	h.HandleFunc("/leaveChatRoom", leaveChatRoom)
	h.HandleFunc("/", noEndpoint)

	//LISTEN AND ERRORS
	err := http.ListenAndServe(":8888", h)
	log.Fatal(err)

	//test := TestNumberDumper(1)
}

//Endpoint nil
func noEndpoint(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintln(w, "You are lost, go home")
}

//CreateChatRoom
func createChatRoom(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println(req.Form)
	fmt.Println(req.Form.Get("q"))
	io.WriteString(w, "hello, world")

}

//ListChatRoom
func listChatRoom(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world")

}

//Join ChatRoom
func jointChatRoom(w http.ResponseWriter, req *http.Request) {

}

//Leave ChatRoom
func leaveChatRoom(w http.ResponseWriter, req *http.Request) {

}

//FOR TESTING
func TestNumberDumper(t *T) {
	n := test(1)
	r, _ := http.NewRequest("GET", "/endpoint1", nil)
	w := httptest.NewRecorder()
	n.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("wrong code returned: %d", w.Code)
	}
	body := w.Body.String()
	if body != fmt.Sprintf("Here is your number: 1\n") {
		t.Fatalf("Wrong body returned: %s", body)
	}
}

func (m test) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Here's your number: %d\n", m)
}
