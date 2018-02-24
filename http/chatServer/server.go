package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"../../lib/ksuid-master"

	"../structs"
)

//type message struct{}
type test int

//type ChatRooms structs.ChatRooms

var chatRooms structs.ChatRooms
var clients structs.Clients

func main() {

	fmt.Println("Starting server...")
	h := http.NewServeMux()

	//DEFINITION OF METHODS
	h.HandleFunc("/createClient", createClient)
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

//Create client
func createClient(w http.ResponseWriter, req *http.Request) {
	fmt.Println("entro aqui")
	//First we get the parameters and parse the  request
	var requestCreateClient structs.RequestCreateClient
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&requestCreateClient)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	var responseCreateClient structs.ResponseCreateClient

	//We execute the creation of client
	for _, client := range clients.Clients {
		if requestCreateClient.Username == client.Username {
			fmt.Println("Client " + requestCreateClient.Username + " already exists")
			responseCreateClient.Status = "Failed, the Client " + requestCreateClient.Username + " already exists"
		}
	}
	currentClient := structs.Client{
		Username: requestCreateClient.Username,
		Id:       ksuid.New().String(),
	}
	clients.Clients = AddClient(currentClient, clients)
	fmt.Println("Client: " + currentClient.Username + " ID: " + currentClient.Id + " created...")
	responseCreateClient.Status = "ok"
	responseCreateClient.Client = currentClient

	//We send the response back
	json.NewEncoder(w).Encode(responseCreateClient)

}

//CreateChatRoom
func createChatRoom(w http.ResponseWriter, req *http.Request) {
	//First We get the parameters
	var requestCreateChatRoom structs.RequestCreateChatRoom
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&requestCreateChatRoom)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//We execute the creation of the chat room
	for _, chatRoom := range chatRooms.Chats {
		if requestCreateChatRoom.NameChatRoom == chatRoom.NameChatRoom {
			fmt.Println("ChatRoom " + requestCreateChatRoom.NameChatRoom + " already exists")
		}
	}
	currentChatRoom := structs.ChatRoom{
		NameChatRoom: requestCreateChatRoom.NameChatRoom,
		Id:           ksuid.New().String(),
		Clients:      structs.Clients{},
	}
	chatRooms.Chats = AddChat(currentChatRoom, chatRooms)
	fmt.Println("ChatRoom: " + currentChatRoom.NameChatRoom + " ID: " + currentChatRoom.Id + " created...")

	//Now we respond to the client
	responseCreateChatRoom := structs.ResponseCreateChatRoom{
		Status:   "ok",
		ChatRoom: currentChatRoom,
	}
	json.NewEncoder(w).Encode(responseCreateChatRoom)

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

//AddChat for append new chats
func AddChat(currentChat structs.ChatRoom, arrayChat structs.ChatRooms) []structs.ChatRoom {
	result := append(arrayChat.Chats, currentChat)
	fmt.Println("Length Chat: ", len(result))
	return result
}

//AddClients for append new Clients
func AddClient(currentClient structs.Client, arrayClient structs.Clients) []structs.Client {
	result := append(arrayClient.Clients, currentClient)
	fmt.Println("Length Clients: ", len(result))
	return result
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
