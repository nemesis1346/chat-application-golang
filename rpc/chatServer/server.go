package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"../../lib/ksuid-master"
	"../structs"
)

//Global variables

func main() {

	chatRooms := new(ChatRooms)
	clients := new(Clients)
	rpc.Register(chatRooms)
	rpc.Register(clients)
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

type ChatRooms structs.ChatRooms
type Clients structs.Clients

//CreateChatRoom...
func (t *ChatRooms) CreateChatRoom(request *structs.RequestCreateChatRoom,
	response *structs.ResponseCreateChatRoom) error {
	currentChatRoom := structs.ChatRoom{
		NameChatRoom: request.NameChatRoom,
		Id:           ksuid.New().String(),
		Clients:      structs.Clients{},
	}
	t.AddChat(currentChatRoom)
	response.Status = "ok"

	return nil
}

//ListChatRoom
func (t *ChatRooms) ListChatRoom(request *structs.RequestListChatRoom,
	response *structs.ResponseListChatRoom) error {
	for _, chatRoom := range t.Chats {
		fmt.Println(chatRoom.NameChatRoom)
	}
	return nil
}
func (t *ChatRooms) ListChatRoomPrint() {
	for _, chatRoom := range t.Chats {
		fmt.Println(chatRoom.NameChatRoom)
	}
}

//CreateClient
func (t *Clients) CreateClient(request *structs.RequestCreateClient,
	response *structs.ResponseCreateClient) error {

	currentClient := structs.Client{
		Username: request.Username,
		Id:       ksuid.New().String(),
	}
	t.AddClient(currentClient)
	fmt.Println("Client: " + currentClient.Username + " ID: " + currentClient.Id + " created...")
	fmt.Println("Length Clients: " + string(len(t.Clients)))
	response.Status = "ok"
	response.Client = currentClient
	return nil
}

//Join ChatRoom
func jointChatRoom(w http.ResponseWriter, req *http.Request) {

}

//Leave ChatRoom
func leaveChatRoom(w http.ResponseWriter, req *http.Request) {

}

//AddChat for append new chats
func (chats *ChatRooms) AddChat(currentChat structs.ChatRoom) []structs.ChatRoom {
	chats.Chats = append(chats.Chats, currentChat)
	return chats.Chats
}

//AddClients for append new Clients
func (clients *Clients) AddClient(currentClient structs.Client) []structs.Client {
	fmt.Println("previous objects: ")
	fmt.Println(clients)
	fmt.Println("current object: ")
	fmt.Println(currentClient)
	clients.Clients = append(clients.Clients, currentClient)
	return clients.Clients
}
