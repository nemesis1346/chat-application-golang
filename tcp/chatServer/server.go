package main

import (
	"encoding/gob"
	"fmt"
	"net"

	"../../lib/ksuid-master"

	"../structs"
)

//type ChatRooms structs.ChatRooms
var chatRooms structs.ChatRooms
var clients structs.Clients

func main() {
	chatRooms = structs.ChatRooms{}
	clients = structs.Clients{}
	clients.Clients = []structs.Client{}

	//We create a listener
	fmt.Println("Starting server...")
	fmt.Println("-----------------------")
	listener, err := net.Listen("tcp", ":12346")
	if err != nil {
		fmt.Println(err)
	}

	//Incomming option messages
	for {
		//We accept the connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
		}
		requestOption := new(structs.OptionMessage)
		//create a decoder object
		gobRequestOption := gob.NewDecoder(conn)
		error := gobRequestOption.Decode(requestOption)
		if error != nil {
			fmt.Println(error)
		}
		fmt.Println("Option: " + requestOption.Option)
		switch requestOption.Option {
		case "1":
			createClient(conn, requestOption)
		case "2":
			createChatRoom(conn)
		case "3":
			listChatRoom(conn)
		case "4":
			joinChatRoom(conn)
		}
	}

}

func createClient(conn net.Conn, requestCreateClient *structs.OptionMessage) {

	flagExists := false
	username := requestCreateClient.Data["Username"]

	//We execute the creation of client
	for _, client := range clients.Clients {
		if username == client.Username {
			flagExists = true
		}
	}
	currentClient := structs.Client{
		Username: username,
		Id:       ksuid.New().String(),
	}
	clients.Clients = AddClient(currentClient, clients)
	fmt.Println("Client: " + currentClient.Username + " ID: " + currentClient.Id + " created...")

	//we send respond back
	mapResCreateClient := make(map[string]string)
	if !flagExists {
		mapResCreateClient["Status"] = "ok"
	} else {
		mapResCreateClient["Status"] = "Client already exists"
	}
	responseCreateClient := structs.OptionMessage{
		Option: "response",
		Data:   mapResCreateClient,
	}
	gobResCreateClient := gob.NewEncoder(conn)
	gobResCreateClient.Encode(responseCreateClient)

}

func createChatRoom(conn net.Conn) {
	fmt.Println("entra aqui")

}
func listChatRoom(conn net.Conn) {

}
func joinChatRoom(conn net.Conn) {

}

// func createClient() {
// 	//request
// 	var requestCreateClient structs.requestCreateClient
// 	gobRequestCreateClient := gob.NewDecoder(requestCreateClient)
// }
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
