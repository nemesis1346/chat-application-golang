package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

//Global variables
var chatRooms []ChatRoom

func main() {
	chatRooms = []ChatRoom{}

	chatRoomGeneric := new(ChatRoom)
	rpc.Register(chatRoomGeneric)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":9999")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)

}

//CreateChatRoom
func (c *ChatRoom) createChatRoom(request *RequestChatRoom, response *ResponseChatRoom) {
	currentChatRoom := &ChatRoom{
		nameChatRoom: request.nameChatRoom,
	}
	chatRooms = append(chatRooms, currentChatRoom)
	response = &ResponseChatRoom{
		status: "ok",
	}

	return
}

type ChatRoom struct {
	nameChatRoom string
}

type RequestChatRoom struct {
	nameChatRoom string
}

type ResponseChatRoom struct {
	status string
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
