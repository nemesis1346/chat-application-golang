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
	//First we check if there is no chats with the request name
	for _, chatRoom := range t.Chats {
		if request.NameChatRoom == chatRoom.NameChatRoom {
			fmt.Println("ChatRoom " + request.NameChatRoom + " already exists")
			response.Status = "Failed, the ChatRoom " + request.NameChatRoom + " already exists"
			return nil
		}
	}
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
	response.ChatRooms.Chats = t.Chats
	response.Status = "ok"
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
	//First we check if the client doesnt exist
	for _, client := range t.Clients {
		if request.Username == client.Username {
			fmt.Println("Client " + request.Username + " already exists")
			response.Status = "Failed, the Client " + request.Username + " already exists"
			return nil
		}
	}
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
			fmt.Print("ChatRoom: " + chatRoom.NameChatRoom + ", Number of Clients: ")
			fmt.Printf("%d\n", len(chatRoom.Clients.Clients))
		}
	}
	return nil
}

func (t *ChatRooms) SendMessages(request *structs.RequestSendMessages,
	response *structs.ResponseSendMessages) {

}

//Leave ChatRoom
func (t *ChatRooms) LeaveChatRoom(request *structs.RequestLeaveChatRoom,
	response *structs.ResponseLeaveChatRoom) error {

	//First we find the chatRoom that name
	counterChat := 0
	for _, chatRoom := range t.Chats {
		if request.ChatRoom.NameChatRoom == chatRoom.NameChatRoom {

			//We find the client we want to erase
			counterClient := 0
			for _, client := range chatRoom.Clients.Clients {
				if request.Client.Username == client.Username {
					fmt.Println("Client erased: " + client.Username)
					DeleteClientFromChatRoom(chatRoom.Clients, counterClient)
					t.Chats[counterChat].Clients = chatRoom.Clients
					fmt.Print("ChatRoom: " + chatRoom.NameChatRoom + ", Number of Clients: ")
					fmt.Printf("%d\n", len(chatRoom.Clients.Clients))
					break
				}
			}
			counterClient++
		}
		counterChat++
	}
	return nil
}

//Get Client
func (t *Clients) GetClient(request *structs.RequestGetClient,
	response *structs.ResponseGetClient) error {
	for _, client := range t.Clients {
		if request.Username == client.Username {
			response.Client = client
			response.Status = "ok"
			return nil
		}
	}
	response.Status = "Not Found"
	return nil
}

//GetChat
func (t *ChatRooms) GetChatRoom(request *structs.RequestGetChatRoom,
	response *structs.ResponseGetChatRoom) error {
	for _, chatRoom := range t.Chats {
		if request.ChatRoomName == chatRoom.NameChatRoom {
			response.ChatRoom = chatRoom
			response.Status = "ok"
			return nil
		}
	}
	response.Status = "Not Found"
	return nil
}

//AddChat for append new chats
func (chats *ChatRooms) AddChat(currentChat structs.ChatRoom) []structs.ChatRoom {
	chats.Chats = append(chats.Chats, currentChat)
	fmt.Println("Length Chat: ")
	fmt.Printf("%d\n", len(chats.Chats))
	return chats.Chats
}

//AddClients for append new Clients
func (clients *Clients) AddClient(currentClient structs.Client) []structs.Client {
	clients.Clients = append(clients.Clients, currentClient)
	fmt.Println("Length Clients: ")
	fmt.Printf("%d\n", len(clients.Clients))
	return clients.Clients
}

//AddClientToChatRoom
func AddClientToChatRoom(currentClient structs.Client, arrayClient structs.Clients) structs.Clients {
	arrayResult := append(arrayClient.Clients, currentClient)
	result := structs.Clients{Clients: arrayResult}
	return result
}

//DeleteClientFromChatRoom
func DeleteClientFromChatRoom(arrayClient structs.Clients, index int) structs.Clients {
	arrayResult := append(arrayClient.Clients[:index], arrayClient.Clients[index+1:]...)
	result := structs.Clients{Clients: arrayResult}
	return result
}
