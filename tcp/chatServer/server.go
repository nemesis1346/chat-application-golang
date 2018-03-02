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
			createChatRoom(conn, requestOption)
		case "3":
			listChatRoom(conn)
		case "4":
			joinChatRoom(conn, requestOption)
		case "5":
			leaveChatRoom(conn, requestOption)
		case "6":
			saveMessage(conn, requestOption)
		case "7":
			getMessages(conn, requestOption)
		}
	}
}

func createClient(conn net.Conn, requestCreateClient *structs.OptionMessage) {
	flagExists := false
	username := requestCreateClient.Data["Username"]

	//We execute the creation of client
	for _, client := range clients.Clients {
		if username == client.Username {
			fmt.Println("Client " + username + " already exists")
			flagExists = true
			break
		}
	}
	mapResCreateClient := make(map[string]string)
	if !flagExists {
		currentClient := structs.Client{
			Username: username,
			Id:       ksuid.New().String(),
		}
		clients.Clients = AddClient(currentClient, clients)
		fmt.Println("Client: " + currentClient.Username + " ID: " + currentClient.Id + " created...")

		mapResCreateClient["Status"] = "ok"
	} else {
		mapResCreateClient["Status"] = "Client already exists"
	}
	//we send back the response to the client
	responseCreateClient := structs.OptionMessage{
		Option: "response",
		Data:   mapResCreateClient,
	}
	gobResCreateClient := gob.NewEncoder(conn)
	gobResCreateClient.Encode(responseCreateClient)
}

func createChatRoom(conn net.Conn, requestCreateChatRoom *structs.OptionMessage) {
	flagExists := false
	fmt.Println("entra aqui")
	nameChatRoom := requestCreateChatRoom.Data["NameChatRoom"]
	//We execute the creation of the chat room
	for _, chatRoom := range chatRooms.Chats {
		if nameChatRoom == chatRoom.NameChatRoom {
			fmt.Println("ChatRoom " + nameChatRoom + " already exists")
			flagExists = true
			break
		}
	}
	mapResCreateChatRoom := make(map[string]string)
	if !flagExists {
		currentChatRoom := structs.ChatRoom{
			NameChatRoom: nameChatRoom,
			Id:           ksuid.New().String(),
			Clients:      structs.Clients{},
		}
		chatRooms.Chats = AddChat(currentChatRoom, chatRooms)
		fmt.Println("ChatRoom: " + nameChatRoom + " ID: " + currentChatRoom.Id + " created...")
		mapResCreateChatRoom["Status"] = "ok"

	} else {
		mapResCreateChatRoom["Status"] = "Client already exists"

	}
	//we send respond back
	responseCreateChatRoom := structs.OptionMessage{
		Option: "response",
		Data:   mapResCreateChatRoom,
	}
	gobResCreateChatRoom := gob.NewEncoder(conn)
	gobResCreateChatRoom.Encode(responseCreateChatRoom)
}

func listChatRoom(conn net.Conn) {
	//we print all the available chats
	for _, chatRoom := range chatRooms.Chats {
		fmt.Print("Name ChatRoom: " + chatRoom.NameChatRoom + " Number of Clients: ")
		fmt.Printf("%d\n", len(chatRoom.Clients.Clients))
		fmt.Println("")
	}
	//we send respond back
	mapResListChatRoom := make(map[string]string)
	mapResListChatRoom["Status"] = "ok"

	//Save all chatrooms in the following map indexes
	//TODO: finish the list

	responseListChatRoom := structs.OptionMessage{
		Option: "response",
		Data:   mapResListChatRoom,
	}
	gobResListChatRoom := gob.NewEncoder(conn)
	gobResListChatRoom.Encode(responseListChatRoom)
}
func joinChatRoom(conn net.Conn, requestJoinChatRoom *structs.OptionMessage) {
	//First We get the parameters
	flagExists := false
	nameChatRoom := requestJoinChatRoom.Data["NameChatRoom"]
	username := requestJoinChatRoom.Data["Username"]

	//GetClient
	currentClient := getClientLocal(username)

	counterChat := 0
	for _, chatRoom := range chatRooms.Chats {
		if nameChatRoom == chatRoom.NameChatRoom {
			arrayClient := AddClientToChatRoom(currentClient, chatRoom.Clients)
			chatRooms.Chats[counterChat].Clients = arrayClient
			fmt.Println("ChatRoom Joined: "+nameChatRoom+", Client: "+currentClient.Username+", Number of Clients: ", len(chatRooms.Chats[counterChat].Clients.Clients))
			flagExists = true
			break
		}
		counterChat++
		//Response failed
	}
	//TODO:OPEN THE LISTENER TO RESPOND IN REAL TIME
	//we send the response back to the client
	mapResJoinChatRoom := make(map[string]string)

	if !flagExists {
		mapResJoinChatRoom["Status"] = "ok"

	} else {
		mapResJoinChatRoom["Status"] = "Client already exists"
	}
	responseJoinChatRoom := structs.OptionMessage{
		Option: "response",
		Data:   mapResJoinChatRoom,
	}
	gobResJoinChatRoom := gob.NewEncoder(conn)
	gobResJoinChatRoom.Encode(responseJoinChatRoom)
}

func leaveChatRoom(conn net.Conn, requestLeaveChatRoom *structs.OptionMessage) {
	//First We get the parameters
	flagExists := false
	nameChatRoom := requestLeaveChatRoom.Data["NameChatRoom"]
	username := requestLeaveChatRoom.Data["Username"]

	//First we find the chatRoom that name
	counterChat := 0
	for _, chatRoom := range chatRooms.Chats {
		fmt.Println("request chat name: " + nameChatRoom + " current chat: " + chatRoom.NameChatRoom)

		if nameChatRoom == chatRoom.NameChatRoom {
			//We find the client we want to erase
			//THE PROBLEM IS THAT THE USER IS NOT IN THE CHATROOM WHEN JOINED
			counterClient := 0
			for _, client := range chatRoom.Clients.Clients {
				if username == client.Username {
					fmt.Println()
					chatRooms.Chats[counterChat].Clients = DeleteClientFromChatRoom(chatRooms.Chats[counterChat].Clients, counterClient)
					//t.Chats[counterChat].Clients = chatRoom.Clients
					fmt.Println("Client  "+client.Username+"left ChatRoom: "+chatRoom.NameChatRoom+", Number of Clients: ", len(chatRooms.Chats[counterChat].Clients.Clients))
					flagExists = true
					break
				}
			}
			counterClient++
		}
		counterChat++
	}
	//we send the response back to the client
	mapResLeaveChatRoom := make(map[string]string)
	if !flagExists {
		mapResLeaveChatRoom["Status"] = "ok"

	} else {
		mapResLeaveChatRoom["Status"] = "Client already exists"
	}
	responseLeaveChatRoom := structs.OptionMessage{
		Option: "response",
		Data:   mapResLeaveChatRoom,
	}
	gobResLeaveChatRoom := gob.NewEncoder(conn)
	gobResLeaveChatRoom.Encode(responseLeaveChatRoom)
}

//Get client
func getClientLocal(username string) structs.Client {
	var result structs.Client
	result = structs.Client{}
	for _, client := range clients.Clients {
		if username == client.Username {
			result = client
		}
	}
	return result
}

//Save message
func saveMessage(conn net.Conn, requestSaveMessage *structs.OptionMessage) {

}

//Get messages
func getMessages(conn net.Conn, requestGetMessages *structs.OptionMessage) {

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
