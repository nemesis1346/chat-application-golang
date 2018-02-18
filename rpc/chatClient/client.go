package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	fmt.Println("Starting client...")
	client, err := rpc.DialHTTP("tcp", ":9999")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Asynchronous call
	request := &RequestChatRoom{
		nameChatRoom: "first try",
	}
	response := &ResponseChatRoom{
		status: "response",
	}
	divCall := client.Go("ChatRoom.createChatRoom", request, response, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	// check errors, print, etc.
	fmt.Println(replyCall)
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
