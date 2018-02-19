package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"reflect"
)

//Global variables

var chatRooms []ChatRoom

func main() {

	chatRoomRA := new(ChatRoomRA)
	rpc.Register(chatRoomRA)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "localhost:1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		//Starting to listen
		fmt.Println("Starting server...")
		fmt.Println("-----------------------")
		http.Serve(l, nil)
	}

}

type ChatRoomRA ChatRoom

type ChatRoom struct {
	nameChatRoom string
}

type RequestChatRoom struct {
	nameChatRoom string
}

type ResponseChatRoom struct {
	status string
}

//CreateChatRoom
func (t *ChatRoomRA) CreateChatRoom(request *RequestChatRoom, response *ResponseChatRoom) error {
	currentChatRoom := &ChatRoom{
		nameChatRoom: request.nameChatRoom,
	}
	fmt.Println(reflect.TypeOf(currentChatRoom))
	//	chatRooms: = append(chatRooms, currentChatRoom)
	response = &ResponseChatRoom{
		status: "ok",
	}

	return nil
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
