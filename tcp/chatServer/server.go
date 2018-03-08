package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"../../lib/ksuid-master"

	"../structs"
)

//type ChatRooms structs.ChatRooms
var chatRooms structs.ChatRooms
var clients structs.Clients
var conn net.Conn //for managing general connection and detect interrumption

func cleanup() {
	fmt.Println("cleanup")
	conn.Close()
}

func main() {

	chatRooms = structs.ChatRooms{}
	clients = structs.Clients{}
	clients.Clients = []structs.Client{}

	//We create a listener
	fmt.Println("Starting server...")
	fmt.Println("-----------------------")
	listener, err := net.Listen("tcp", ":12346")
	if err != nil {
		fmt.Printf("[Listener]\t", err)
	}
	for {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			cleanup()
			os.Exit(1)
		}()
		//We accept the connection
		conn, err = listener.Accept()
		if err != nil {
			fmt.Printf("[Connection]\t", err)
			conn.Close()
		}
		go handleClient(conn)
	}

}
func handleClient(conn net.Conn) {
	//Incomming option messages
	for {
		requestOption := new(structs.OptionMessage)
		//create a decoder object
		gobRequestOption := gob.NewDecoder(conn)
		err := gobRequestOption.Decode(requestOption)
		if err != nil {
			fmt.Printf("[RequestOption]\t", err)
			break
		}
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
			getPreviousMessages(conn, requestOption)
		case "6":
			leaveChatRoom(conn, requestOption)
		case "7":
			saveMessage(conn, requestOption)
		case "8":
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
	mapResListChatRoom := make(map[string]string)

	//we print all the available chats
	if len(chatRooms.Chats) > 0 {
		mapResListChatRoom["Status"] = "ok"
		//we send respond back
		for _, chatRoom := range chatRooms.Chats {
			stringLength := strconv.Itoa(len(chatRoom.Clients.Clients))
			mapResListChatRoom[chatRoom.NameChatRoom] = stringLength
			fmt.Print("Name ChatRoom: " + chatRoom.NameChatRoom + " Number of Clients: " + stringLength)
			fmt.Println("")
		}
	} else {
		mapResListChatRoom["Status"] = "There is not chatRooms available"
	}

	//We send back the response to the client

	responseListChatRoom := structs.OptionMessage{
		Option: "response",
		Data:   mapResListChatRoom,
	}
	gobResListChatRoom := gob.NewEncoder(conn)
	gobResListChatRoom.Encode(responseListChatRoom)
}

func getPreviousMessages(conn net.Conn, requestGetPreviousMessages *structs.OptionMessage) {
	//First we get the parameters
	nameChatRoom := requestGetPreviousMessages.Data["NameChatRoom"]

	//We create the response
	resGetPreviousMessages := make(map[string]string)

	//Now we get the previous messages
	var arrayMessages structs.Messages

	for _, chatRoom := range chatRooms.Chats {
		if chatRoom.NameChatRoom == nameChatRoom {
			arrayMessages = chatRoom.Messages
			break
		}
	}
	if len(arrayMessages.Messages) > 0 {
		for _, message := range arrayMessages.Messages {
			resGetPreviousMessages[message.Username+":"+message.Content+" Time: "] = message.Time.Format(time.RFC3339)
		}
		resGetPreviousMessages["Status"] = "ok"
	} else {
		resGetPreviousMessages["Status"] = "There are not previous messages..."
	}
	//We send the response back
	resPreviousMessages := structs.OptionMessage{
		Option: "response",
		Data:   resGetPreviousMessages,
	}
	gobResPreviousMessages := gob.NewEncoder(conn)
	gobResPreviousMessages.Encode(resPreviousMessages)
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

	if flagExists {
		mapResJoinChatRoom["Status"] = "ok"

	} else {
		mapResJoinChatRoom["Status"] = "There is no chat with that name"
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
		mapResLeaveChatRoom["Status"] = "There was some error leaving chat room"
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
	//First we get the parameters
	content := requestSaveMessage.Data["Content"]
	username := requestSaveMessage.Data["Username"]
	nameChatRoom := requestSaveMessage.Data["NameChatRoom"]
	timeRequestString := requestSaveMessage.Data["TimeRequest"]
	layout := time.RFC3339
	timeRequest, _ := time.Parse(layout, timeRequestString)

	//First create message instance
	currentMessage := structs.Message{
		Id:           ksuid.New().String(),
		Content:      content,
		Username:     username,
		NameChatRoom: nameChatRoom,
		Time:         time.Now(),
	}

	//Initial value of response
	mapResSaveMessage := make(map[string]string)
	mapResSaveMessage["Status"] = "There was a problem in saving the message"

	counterChat := 0
	//Find the chat with the name from the request
	for _, chatRoom := range chatRooms.Chats {
		if chatRoom.NameChatRoom == nameChatRoom {
			arrayMessages := AddMessagesInChatRoom(currentMessage, chatRoom.Messages)
			chatRooms.Chats[counterChat].Messages = arrayMessages
			fmt.Println("length: " + strconv.Itoa(len(chatRooms.Chats[counterChat].Messages.Messages)))
			//We save the messages in the response
			fmt.Println(username + " :" + content + " Time: " + currentMessage.Time.Format(time.RFC3339))

			mapResSaveMessage["Status"] = "ok"
			//mapResSaveMessage["TimeLast"] = currentMessage.Time.Format(time.RFC3339)
			mapResSaveMessage["TimeLast"] = timeRequest.Format(time.RFC3339)
			break
		}
		counterChat++
	}
	//we send the response back to the client

	responseSaveMessage := structs.OptionMessage{
		Option: "response",
		Data:   mapResSaveMessage,
	}
	gobResSaveMessage := gob.NewEncoder(conn)
	gobResSaveMessage.Encode(responseSaveMessage)
}

//Get messages
func getMessages(conn net.Conn, requestGetMessages *structs.OptionMessage) {
	//First we get the parameters
	timeRequestString := requestGetMessages.Data["Time"]
	//fmt.Println("Request time : " + timeRequestString)
	nameChatRoom := requestGetMessages.Data["NameChatRoom"]

	layout := time.RFC3339
	timestamp, err := time.Parse(layout, timeRequestString)
	if err != nil {
		fmt.Printf("[GetMessages]\t", err)
	}

	//Variables to respond
	mapResGetMessages := make(map[string]string)
	mapResGetMessages["Status"] = "We could not get the messages"
	//Now we get the messages
	counterChat := 0

	for _, chatRoom := range chatRooms.Chats {
		if chatRoom.NameChatRoom == nameChatRoom {
			if len(chatRooms.Chats[counterChat].Messages.Messages) > 0 {
				//fmt.Println("listening msn: " + strconv.Itoa(len(chatRooms.Chats[counterChat].Messages.Messages)))
				for _, message := range chatRooms.Chats[counterChat].Messages.Messages {
					if message.Time.After(timestamp) {
						mapResGetMessages[message.Username+" : "+message.Content+" Time:"] = string(message.Time.Format(time.RFC3339))
					}

				}
				mapResGetMessages["Status"] = "ok"
			} else {
				//There is not any message between that period of time
				mapResGetMessages["Status"] = "There is no previous messages"
			}
			break
		}
		counterChat++
	}
	//we send the response back to the client

	responseGetMessages := structs.OptionMessage{
		Option: "response",
		Data:   mapResGetMessages,
	}
	gobResGetMessages := gob.NewEncoder(conn)
	gobResGetMessages.Encode(responseGetMessages)
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
