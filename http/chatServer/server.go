package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"../../lib/ksuid-master"

	"../structs"
)

//type ChatRooms structs.ChatRooms
var chatRooms structs.ChatRooms
var clients structs.Clients

//var bundleMessages structs.Messages

func main() {

	fmt.Println("Starting server...")
	h := http.NewServeMux()

	//DEFINITION OF METHODS
	h.HandleFunc("/createClient", createClient)
	h.HandleFunc("/createChatRoom", createChatRoom)
	h.HandleFunc("/listChatRoom", listChatRoom)
	h.HandleFunc("/joinChatRoom", jointChatRoom)
	h.HandleFunc("/leaveChatRoom", leaveChatRoom)
	h.HandleFunc("/saveMessage", saveMessage)
	h.HandleFunc("/getChatRoom", getChatRoom)
	h.HandleFunc("/getPreviousMessages", getPreviousMessages)
	h.HandleFunc("/getMessages", getMessages)
	h.HandleFunc("/", noEndpoint)

	//LISTEN AND ERRORS
	err := http.ListenAndServe(":8888", h)
	log.Fatal(err)
}

//Endpoint nil
func noEndpoint(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintln(w, "You are lost, go home")
}

//Get Messages
func getMessages(w http.ResponseWriter, req *http.Request) {
	//First We get the parameters
	var requestGetMessages structs.RequestGetMessages
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&requestGetMessages)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var responseGetMessages structs.ResponseGetMessages
	responseGetMessages.Status = "There is no previous messages"

	//Now we get the messages
	counterChat := 0

	for _, chatRoom := range chatRooms.Chats {
		if chatRoom.NameChatRoom == requestGetMessages.ChatRoom.NameChatRoom {
			if len(chatRooms.Chats[counterChat].Messages.Messages) > 0 {
				arrayMessages := structs.Messages{}
				for _, message := range chatRooms.Chats[counterChat].Messages.Messages {
					if message.Time.After(requestGetMessages.Time) {
						arrayMessages = AddMessagesInChatRoom(message, arrayMessages)
					}
				}
				responseGetMessages.Messages = arrayMessages
				responseGetMessages.Status = "ok"
			}
			break
		}
		counterChat++
	}
	//we send back the response to the client
	json.NewEncoder(w).Encode(responseGetMessages)
}

//Create client
func createClient(w http.ResponseWriter, req *http.Request) {
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
	//First We get the parameters
	var requestListChatRoom structs.RequestListChatRoom
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&requestListChatRoom)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//We execute the creation of the chat room

	for _, chatRoom := range chatRooms.Chats {
		fmt.Print("Name ChatRoom: " + chatRoom.NameChatRoom + " Number of Clients: ")
		fmt.Printf("%d\n", len(chatRoom.Clients.Clients))
		fmt.Println("")
	}

	//we respond to the clinet
	responseListChatRoom := structs.ResponseListChatRoom{
		ChatRooms: chatRooms,
		Status:    "ok",
	}
	json.NewEncoder(w).Encode(responseListChatRoom)

}

//Join ChatRoom
func jointChatRoom(w http.ResponseWriter, req *http.Request) {
	//First We get the parameters
	var requestJoinChatRoom structs.RequestJoinChatRoom
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&requestJoinChatRoom)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//Now we join the client
	var responseJoinChatRoom structs.ResponseJoinChatRoom

	//Initial value of answer
	responseJoinChatRoom.Status = "There is no chatRoom with that name"

	counterChat := 0
	for _, chatRoom := range chatRooms.Chats {
		if chatRoom.NameChatRoom == requestJoinChatRoom.ChatRoom.NameChatRoom {
			arrayClient := AddClientToChatRoom(requestJoinChatRoom.Client, chatRoom.Clients)
			chatRooms.Chats[counterChat].Clients = arrayClient
			fmt.Println("ChatRoom Joined: "+chatRoom.NameChatRoom+", Client: "+requestJoinChatRoom.Client.Username+", Number of Clients: ", len(chatRooms.Chats[counterChat].Clients.Clients))
			responseJoinChatRoom.Status = "ok"
			break
		}
		counterChat++
		//Response failed
	}
	//we send the response back to the client
	json.NewEncoder(w).Encode(responseJoinChatRoom)

}

//Save message
func saveMessage(w http.ResponseWriter, req *http.Request) {
	//First We get the parameters
	var requestSaveMessage structs.RequestSaveMessage
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&requestSaveMessage)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Println(requestSaveMessage.Client.Username + ": " + requestSaveMessage.Content + "  Time: " + requestSaveMessage.Time.Format(time.RFC3339))

	//First create message instance
	currentMessage := structs.Message{
		Id:           ksuid.New().String(),
		Content:      requestSaveMessage.Content,
		Username:     requestSaveMessage.Client.Username,
		NameChatRoom: requestSaveMessage.ChatRoom.NameChatRoom,
		Time:         time.Now(),
	}

	var responseSaveMessage structs.ResponseSaveMessage
	//Initial value of response
	responseSaveMessage.Status = "There is no chat room with that name"

	counterChat := 0
	//Find the chat with the name from the request
	for _, chatRoom := range chatRooms.Chats {
		if chatRoom.NameChatRoom == requestSaveMessage.ChatRoom.NameChatRoom {
			arrayMessages := AddMessagesInChatRoom(currentMessage, chatRoom.Messages)
			chatRooms.Chats[counterChat].Messages = arrayMessages
			//We save the messages in the response
			responseSaveMessage.Messages = arrayMessages
			responseSaveMessage.Status = "ok"
			responseSaveMessage.Time = requestSaveMessage.Time //TODO: maybe i dont need this
			break
		}
		counterChat++
	}
	//we send the response back to the client
	json.NewEncoder(w).Encode(responseSaveMessage)
}

//Leave ChatRoom
func leaveChatRoom(w http.ResponseWriter, req *http.Request) {

	//First We get the parameters
	var requestLeaveChatRoom structs.RequestLeaveChatRoom
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&requestLeaveChatRoom)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var responseLeaveChatRoom structs.ResponseLeaveChatRoom
	//initial status
	responseLeaveChatRoom.Status = "Problem in leaving chatroom"

	//First we find the chatRoom that name
	counterChat := 0
	for _, chatRoom := range chatRooms.Chats {
		fmt.Println("request chat name: " + requestLeaveChatRoom.ChatRoom.NameChatRoom + " current chat: " + chatRoom.NameChatRoom)

		if requestLeaveChatRoom.ChatRoom.NameChatRoom == chatRoom.NameChatRoom {
			//We find the client we want to erase
			//THE PROBLEM IS THAT THE USER IS NOT IN THE CHATROOM WHEN JOINED
			counterClient := 0
			for _, client := range chatRoom.Clients.Clients {
				if requestLeaveChatRoom.Client.Username == client.Username {
					fmt.Println()
					chatRooms.Chats[counterChat].Clients = DeleteClientFromChatRoom(chatRooms.Chats[counterChat].Clients, counterClient)
					//t.Chats[counterChat].Clients = chatRoom.Clients
					fmt.Println("Client  "+client.Username+"left ChatRoom: "+chatRoom.NameChatRoom+", Number of Clients: ", len(chatRooms.Chats[counterChat].Clients.Clients))
					responseLeaveChatRoom.Status = "ok"
					break
				}
			}
			counterClient++
		}
		counterChat++
	}
	//we send the response back to the client
	json.NewEncoder(w).Encode(responseLeaveChatRoom)
}

//Get previous messages
func getPreviousMessages(w http.ResponseWriter, req *http.Request) {
	//First We get the parameters
	var requestGetPreviousMessages structs.RequestGetPreviousMessages
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&requestGetPreviousMessages)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var responseGetPreviousMessages structs.ResponseGetPreviousMessages
	//Now we get the previous messages
	for _, chatRoom := range chatRooms.Chats {
		if chatRoom.NameChatRoom == requestGetPreviousMessages.ChatRoom.NameChatRoom {
			responseGetPreviousMessages.Messages = chatRoom.Messages
			responseGetPreviousMessages.Status = "ok"
			break
		}
		responseGetPreviousMessages.Status = "There is no previous messages"
	}
	//we send back the response to the client
	json.NewEncoder(w).Encode(responseGetPreviousMessages)

}

//Get ChatRoom
func getChatRoom(w http.ResponseWriter, req *http.Request) {
	//First We get the parameters
	var requestGetChatRoom structs.RequestGetChatRoom
	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(req.Body).Decode(&requestGetChatRoom)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var responseGetChatRoom structs.ResponseGetChatRoom

	for _, chatRoom := range chatRooms.Chats {
		if requestGetChatRoom.ChatRoomName == chatRoom.NameChatRoom {
			responseGetChatRoom.ChatRoom = chatRoom
			responseGetChatRoom.Status = "ok"
			break
		}
		responseGetChatRoom.Status = "Not Found"
	}
	//we send the response back to the client
	json.NewEncoder(w).Encode(responseGetChatRoom)
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
