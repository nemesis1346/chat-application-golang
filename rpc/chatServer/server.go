package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"../../lib/ksuid-master"
	"../structs"
)

//Global variables
type ChatRooms structs.ChatRooms
type Clients structs.Clients
type Client structs.Client
type ChatRoom structs.ChatRoom

func main() {

	//We start server
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
	counterChat := 0
	for _, chatRoom := range t.Chats {
		if chatRoom.NameChatRoom == request.ChatRoom.NameChatRoom {
			arrayClient := AddClientToChatRoom(request.Client, chatRoom.Clients)
			t.Chats[counterChat].Clients = arrayClient
			fmt.Println("ChatRoom Joined: "+chatRoom.NameChatRoom+", Client: "+request.Client.Username+", Number of Clients: ", len(t.Chats[counterChat].Clients.Clients))
			response.Status = "ok"
			return nil
		}
		counterChat++
	}
	response.Status = "There is no chatRoom with that name"
	return nil
}

//Save current message
func (t *ChatRooms) SaveMessage(request *structs.RequestSaveMessage,
	response *structs.ResponseSaveMessage) error {
	fmt.Println(request.Client.Username + ": " + request.Content + "  Time: " + request.Time.Format(time.RFC3339))

	//First create message instance
	currentMessage := structs.Message{
		Id:           ksuid.New().String(),
		Content:      request.Content,
		Username:     request.Client.Username,
		NameChatRoom: request.ChatRoom.NameChatRoom,
		Time:         request.Time,
	}
	counterChat := 0
	//Find the chat with the name from the request
	for _, chatRoom := range t.Chats {
		if chatRoom.NameChatRoom == request.ChatRoom.NameChatRoom {
			arrayMessages := AddMessagesInChatRoom(currentMessage, chatRoom.Messages)
			t.Chats[counterChat].Messages = arrayMessages
			//We save the messages in the response
			response.Messages = arrayMessages
			response.Status = "ok"
			response.Time = request.Time
			return nil
		}
		counterChat++
	}
	return nil
}

//Get list of previous messages
func (t *ChatRooms) GetPreviousMessages(request *structs.RequestGetPreviousMessages,
	response *structs.ResponseGetPreviousMessages) error {
	for _, chatRoom := range t.Chats {
		if chatRoom.NameChatRoom == request.ChatRoom.NameChatRoom {
			response.Messages = chatRoom.Messages
			response.Status = "ok"
			return nil
		}
	}
	return nil
}

//Leave ChatRoom
func (t *ChatRooms) LeaveChatRoom(request *structs.RequestLeaveChatRoom,
	response *structs.ResponseLeaveChatRoom) error {
	//First we find the chatRoom that name
	counterChat := 0
	for _, chatRoom := range t.Chats {
		fmt.Println("request chat name: " + request.ChatRoom.NameChatRoom + " current chat: " + chatRoom.NameChatRoom)

		if request.ChatRoom.NameChatRoom == chatRoom.NameChatRoom {
			//We find the client we want to erase
			//THE PROBLEM IS THAT THE USER IS NOT IN THE CHATROOM WHEN JOINED
			counterClient := 0
			for _, client := range chatRoom.Clients.Clients {
				if request.Client.Username == client.Username {
					fmt.Println("Client erased: " + client.Username)
					t.Chats[counterChat].Clients = DeleteClientFromChatRoom(t.Chats[counterChat].Clients, counterClient)
					//t.Chats[counterChat].Clients = chatRoom.Clients
					fmt.Println("ChatRoom: "+chatRoom.NameChatRoom+", Number of Clients: ", len(t.Chats[counterChat].Clients.Clients))
					response.Status = "ok"
					return nil
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
	fmt.Println("Length Chat: ", len(chats.Chats))
	return chats.Chats
}

//AddClients for append new Clients
func (clients *Clients) AddClient(currentClient structs.Client) []structs.Client {
	clients.Clients = append(clients.Clients, currentClient)
	fmt.Println("Length Clients: ", len(clients.Clients))
	return clients.Clients
}

//AddClientToChatRoom
func AddClientToChatRoom(currentClient structs.Client, arrayClient structs.Clients) structs.Clients {
	arrayResult := append(arrayClient.Clients, currentClient)
	result := structs.Clients{Clients: arrayResult}
	return result
}

//AddMessages
func AddMessagesInChatRoom(currentMessage structs.Message, arrayMessage structs.Messages) structs.Messages {
	arrayResult := append(arrayMessage.Messages, currentMessage)
	result := structs.Messages{Messages: arrayResult}
	return result
}

//DeleteClientFromChatRoom
func DeleteClientFromChatRoom(arrayClient structs.Clients, index int) structs.Clients {
	arrayResult := append(arrayClient.Clients[:index], arrayClient.Clients[index+1:]...)
	result := structs.Clients{Clients: arrayResult}
	return result
}
