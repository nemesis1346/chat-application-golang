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

	//We create a listener
	fmt.Println("Starting server...")
	fmt.Println("-----------------------")
	listener, err := net.Listen("tcp", ":12346/createClient")
	if err != nil {
		fmt.Println(err)
	}

	//now we create new client

	//We create a buffer
	// tmp := make([]byte, 1024)
	// fmt.Println(tmp)

	// _, err = conn.Read(tmp)
	// fmt.Println(tmp)

	//We take off extra zeros
	//tmp = bytes.Trim(tmp, "\x00")

	//fmt.Println(tmp)

	//Temporal buffer for incomming messages
	//tmpBuffer := bytes.NewBuffer(tmp)

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
		fmt.Println(error)
		fmt.Println(requestOption.Option)
		switch requestOption.Option {
		case "1":
			createChatRoom(conn)
		case "2":
			listChatRoom(conn)
		case "3":
			joinChatRoom(conn)
		}
		conn.Close()
	}
}

func createClient(conn net.Conn, requestCreateClient structs.Client) {
	var responseCreateClient structs.ResponseCreateClient

	//We execute the creation of client
	for _, client := range clients.Clients {
		if requestCreateClient.Username == client.Username {
			fmt.Println("Client " + requestCreateClient.Username + " already exists")
			responseCreateClient.Status = "Failed, the Client " + requestCreateClient.Username + " already exists"
		}
		currentClient := structs.Client{
			Username: requestCreateClient.Username,
			Id:       ksuid.New().String(),
		}
		clients.Clients = AddClient(currentClient, clients)
		fmt.Println("Client: " + currentClient.Username + " ID: " + currentClient.Id + " created...")
		responseCreateClient.Status = "ok"
		responseCreateClient.Client = currentClient
	}

}

func createChatRoom(conn net.Conn) {

	fmt.Print("create chat room")
	conn.Close()

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
