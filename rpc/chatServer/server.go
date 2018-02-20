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
	clients.Clients = []structs.Client{}
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
type Client structs.Client
type ChatRoom structs.ChatRoom

//CreateChatRoom...
func (t *ChatRooms) CreateChatRoom(request *structs.RequestCreateChatRoom,
	response *structs.ResponseCreateChatRoom) error {
	currentChatRoom := structs.ChatRoom{
		NameChatRoom: request.NameChatRoom,
		Id:           ksuid.New().String(),
		Clients:      structs.Clients{},
	}
	t.AddChat(currentChatRoom)
	fmt.Println("ChatRoom: " + currentChatRoom.NameChatRoom + " ID: " + currentChatRoom.Id + " created...")
	response.Status = "ok"
	response.ChatRoom = currentChatRoom
	return nil
}

//ListChatRoom
func (t *ChatRooms) ListChatRoom(request *structs.RequestListChatRoom,
	response *structs.ResponseListChatRoom) error {
	for _, chatRoom := range t.Chats {
		fmt.Print("Name ChatRoom: " + chatRoom.NameChatRoom + " Number of Clients: ")
		fmt.Printf("%d\n", len(chatRoom.Clients.Clients))
		fmt.Println("")
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
	response.Status = "ok"
	response.Client = currentClient
	return nil
}

//Join ChatRoom
func (t *ChatRooms) JoinChatRoom(request *structs.RequestJoinChatRoom,
	response *structs.ResponseJoinChatRoom) error {
	for _, chatRoom := range t.Chats {
		if chatRoom.NameChatRoom == request.ChatRoom.NameChatRoom {
			arrayClient := AddClientToChatRoom(request.Client, chatRoom.Clients)
			chatRoom.Clients = arrayClient
			fmt.Print("ChatRoom: " + chatRoom.NameChatRoom + " Number of Clients: ")
			fmt.Printf("%d\n", len(chatRoom.Clients.Clients))
		}
	}
	return nil
}

//Leave ChatRoom
func (t *ChatRooms) leaveChatRoom(request *structs.RequestLeaveChatRoom,
	response *structs.ResponseLeaveChatRoom) error {
	return nil
}

//AddChat for append new chats
func (chats *ChatRooms) AddChat(currentChat structs.ChatRoom) []structs.ChatRoom {
	chats.Chats = append(chats.Chats, currentChat)
	fmt.Println("Length: ")
	fmt.Printf("%d\n", len(chats.Chats))
	return chats.Chats
}

//AddClients for append new Clients
func (clients *Clients) AddClient(currentClient structs.Client) []structs.Client {
	// fmt.Println("previous objects: ")
	// fmt.Println(clients)
	// fmt.Println("current object: ")
	// fmt.Println(currentClient)
	clients.Clients = append(clients.Clients, currentClient)
	fmt.Println("Length: ")
	fmt.Printf("%d\n", len(clients.Clients))
	return clients.Clients
}
func AddClientToChatRoom(currentClient structs.Client, arrayClient structs.Clients) structs.Clients {
	arrayResult := append(arrayClient.Clients, currentClient)
	result := structs.Clients{Clients: arrayResult}
	return result
}
