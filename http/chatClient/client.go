package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"../structs"
)

//var time
var currentTime time.Time

func main() {

	fmt.Println("Starting client...")

	//Introduce credentials
	fmt.Println("Enter a username: ")
	username := bufio.NewReader(os.Stdin)

	inputUserObject, err := username.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	inputUserObject = inputUserObject[:len(inputUserObject)-1]

	//Now we create the client
	requestCreateClient := structs.RequestCreateClient{
		Username: string(inputUserObject)}

	bCreateClient := new(bytes.Buffer)
	json.NewEncoder(bCreateClient).Encode(requestCreateClient)

	resCreateClient, errorResCreateClient := http.Post("http://localhost:8888/createClient", "application/json; charset=utf-8", bCreateClient)
	if errorResCreateClient != nil {
		fmt.Println(errorResCreateClient)
	}
	var bodyCreateClient structs.ResponseCreateClient
	json.NewDecoder(resCreateClient.Body).Decode(&bodyCreateClient)

	//If the call was successful and the client is created, we start chosing options
	if bodyCreateClient.Status == "ok" {
		fmt.Println()
		fmt.Println("----User: " + bodyCreateClient.Client.Username + " created!!")
		fmt.Println()
		currentGlobalUser := bodyCreateClient.Client
		fmt.Println("Current User: " + currentGlobalUser.Username + " Id: " + currentGlobalUser.Id)

		for {
			//Introduce options
			option := bufio.NewReader(os.Stdin)

			//Interface for options of the client
			fmt.Println("Choose an action:")
			fmt.Println("1.Create a chatroom")
			fmt.Println("2.List all existing chatrooms ")
			fmt.Println("3.Join a chatroom ")
			fmt.Println("4.Leave a chatroom \n")

			fmt.Print("Choose option: ")
			//reading the input
			input, _, err := option.ReadRune()
			if err != nil {
				fmt.Println(err)
			}

			//OPTIONS
			switch input {
			case '1':
				createChatRoom()
			case '2':
				listChatRoom(currentGlobalUser)
			case '3':
				joinChatRoom(currentGlobalUser)
			}
		}
	} else {
		fmt.Println("There was some error in client creation")
	}

}

//listen to incomming messages
func listenMessages(currentUser structs.Client, chatRoom structs.ChatRoom) {
	for {
		requestGetMessages := structs.RequestGetMessages{
			ChatRoom: chatRoom,
			Client:   currentUser,
			Time:     currentTime,
		}

		bGetMessages := new(bytes.Buffer)
		json.NewEncoder(bGetMessages).Encode(requestGetMessages)
		resGetMessages, _ := http.Post("http://localhost:8888/getMessages", "application/json; charset=utf-8", bGetMessages)

		var bodyGetMessages structs.ResponseGetMessages
		json.NewDecoder(resGetMessages.Body).Decode(&bodyGetMessages)

		if len(bodyGetMessages.Messages.Messages) > 0 {
			for _, messageResult := range bodyGetMessages.Messages.Messages {
				fmt.Println(messageResult.Username + ": " + messageResult.Content + "  Time: " + messageResult.Time.Format(time.RFC3339))

			}
			currentTime = time.Now()
		}
		time.Sleep(time.Second / 2)
	}
}

//Create chat room
func createChatRoom() {
	//Input of chatname
	fmt.Print("Choose a name for the chatRoom: ")

	reader := bufio.NewReader(os.Stdin)
	chatName, _ := reader.ReadString('\n')
	chatName = chatName[:len(chatName)-1]

	requestCreateChatRoom := structs.RequestCreateChatRoom{
		NameChatRoom: string(chatName)}

	bCreateChatRoom := new(bytes.Buffer)
	json.NewEncoder(bCreateChatRoom).Encode(requestCreateChatRoom)

	resCreateChatRoom, _ := http.Post("http://localhost:8888/createChatRoom", "application/json; charset=utf-8", bCreateChatRoom)

	var bodyCreateChatRoom structs.ResponseCreateChatRoom
	json.NewDecoder(resCreateChatRoom.Body).Decode(&bodyCreateChatRoom)

	if bodyCreateChatRoom.Status == "ok" {
		fmt.Println()
		fmt.Println("----ChatCreated: " + bodyCreateChatRoom.Status)
		fmt.Println("NameChatRoom: " + bodyCreateChatRoom.ChatRoom.NameChatRoom)
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println("----Error: " + bodyCreateChatRoom.Status)
		fmt.Println()
	}
}

func listChatRoom(currentUser structs.Client) {
	requestListChatRoom := structs.RequestListChatRoom{
		Username: string(currentUser.Username)}

	//We execute the call to the server
	bListChatRoom := new(bytes.Buffer)
	json.NewEncoder(bListChatRoom).Encode(requestListChatRoom)
	resListChatRoom, _ := http.Post("http://localhost:8888/listChatRoom", "application/json; charset=utf-8", bListChatRoom)
	var bodyListChatRoom structs.ResponseListChatRoom
	json.NewDecoder(resListChatRoom.Body).Decode(&bodyListChatRoom)

	fmt.Println()
	fmt.Println("Chats available .....")
	fmt.Println()

	resultArray := bodyListChatRoom.ChatRooms.Chats

	//Now we print the result
	if len(resultArray) > 0 {
		for _, chatRoom := range resultArray {
			fmt.Print("Name ChatRoom: " + chatRoom.NameChatRoom + ", Number of Clients: ")
			fmt.Printf("%d\n", len(chatRoom.Clients.Clients))
			fmt.Println("")
		}
	} else {
		fmt.Println("There is no chats available")
		fmt.Println()
	}

}
func joinChatRoom(currentUser structs.Client) {
	//First we show the chats available
	listChatRoom(currentUser)

	//Input of chatname
	fmt.Print("Choose a name for joining chatRoom: ")

	reader := bufio.NewReader(os.Stdin)
	chatName, _ := reader.ReadString('\n')
	chatName = chatName[:len(chatName)-1]

	//First we get the chatRoom struct
	requestGetChatRoom := structs.RequestGetChatRoom{
		ChatRoomName: chatName,
	}
	//We execute the call to the server to get the chatRoom
	bGetChatRoom := new(bytes.Buffer)
	json.NewEncoder(bGetChatRoom).Encode(requestGetChatRoom)
	resGetChatRoom, _ := http.Post("http://localhost:8888/getChatRoom", "application/json; charset=utf-8", bGetChatRoom)
	var bodyGetChatRoom structs.ResponseGetChatRoom
	json.NewDecoder(resGetChatRoom.Body).Decode(&bodyGetChatRoom)

	if bodyGetChatRoom.Status == "ok" {
		//Now we submit request for joining chatRoom
		requestJoinChatRoom := structs.RequestJoinChatRoom{
			Client:   currentUser,
			ChatRoom: bodyGetChatRoom.ChatRoom,
		}
		//We make the call to join to the chat room
		bJoinChatRoom := new(bytes.Buffer)
		json.NewEncoder(bJoinChatRoom).Encode(requestJoinChatRoom)
		resJoinChatRoom, _ := http.Post("http://localhost:8888/joinChatRoom", "application/json; charset=utf-8", bJoinChatRoom)
		var bodyJoinChatRoom structs.ResponseJoinChatRoom
		json.NewDecoder(resJoinChatRoom.Body).Decode(&bodyJoinChatRoom)

		//we evaluate any error
		if bodyJoinChatRoom.Status == "ok" {

			//we print the client joining
			fmt.Println()
			fmt.Println("Client " + currentUser.Username + " joined to " + requestJoinChatRoom.ChatRoom.NameChatRoom)
			fmt.Println()

			//Now we get the rest of the messages
			fmt.Println("Previous messages from " + requestJoinChatRoom.ChatRoom.NameChatRoom + " ...")
			requestGetPreviousMessages := structs.RequestGetPreviousMessages{
				ChatRoom: requestJoinChatRoom.ChatRoom,
				Client:   currentUser,
			}
			//We make the call for the previous messages
			bPreviousMessages := new(bytes.Buffer)
			json.NewEncoder(bPreviousMessages).Encode(requestGetPreviousMessages)
			resPreviousMessages, _ := http.Post("http://localhost:8888/getPreviousMessages", "application/json; charset=utf-8", bPreviousMessages)
			var bodyPreviousMessages structs.ResponseGetPreviousMessages
			json.NewDecoder(resPreviousMessages.Body).Decode(&bodyPreviousMessages)

			previousMessages := bodyPreviousMessages.Messages.Messages

			for _, message := range previousMessages {
				fmt.Println(message.Username + ": " + message.Content)
			}

			//current message time TODO: the time should be the last time of the message
			currentTime = time.Now()
			//Now we start chating
			fmt.Println("Start chating.....")
			fmt.Println()

			go listenMessages(currentUser, bodyGetChatRoom.ChatRoom) //Start listening messages

			for {
				readerMessage := bufio.NewReader(os.Stdin)
				messageContent, _ := readerMessage.ReadString('\n')
				messageContent = messageContent[:len(messageContent)]
				//First we detect the command exit
				if strings.TrimRight(messageContent, "\n") == "exit" {
					leaveChatRoom(currentUser, bodyGetChatRoom.ChatRoom.NameChatRoom)
					break
				} else {
					//We submit request for the message

					requestSaveMessage := structs.RequestSaveMessage{
						Client:   currentUser,
						Content:  messageContent,
						ChatRoom: bodyGetChatRoom.ChatRoom,
						Time:     currentTime,
					}
					//we make the call for saving messages
					bSaveMessage := new(bytes.Buffer)
					json.NewEncoder(bSaveMessage).Encode(requestSaveMessage)
					resSaveMessage, _ := http.Post("http://localhost:8888/saveMessage", "application/json; charset=utf-8", bSaveMessage)
					var bodySaveMessage structs.ResponseSaveMessage
					json.NewDecoder(resSaveMessage.Body).Decode(&bodySaveMessage)

					if bodySaveMessage.Status != "ok" {
						fmt.Println("There was an error in saving message")
					}
					messagesResponse := bodySaveMessage.Messages.Messages
					currentTime = messagesResponse[len(messagesResponse)-1].Time

				}
			}
		} else {
			fmt.Println("There was some error in join chat room")
		}
	} else {
		fmt.Println("There was some error in get chat room")
	}
}

func leaveChatRoom(currentUser structs.Client, chatName string) {
	//First we get the chatRoom struct
	requestGetChatRoom := structs.RequestGetChatRoom{
		ChatRoomName: chatName,
	}

	//we make the call for getting the chat room
	bGetChatRoom := new(bytes.Buffer)
	json.NewEncoder(bGetChatRoom).Encode(requestGetChatRoom)
	resGetChatRoom, _ := http.Post("http://localhost:8888/getChatRoom", "application/json; charset=utf-8", bGetChatRoom)
	var bodyGetChatRoom structs.ResponseGetChatRoom
	json.NewDecoder(resGetChatRoom.Body).Decode(&bodyGetChatRoom)

	if bodyGetChatRoom.Status == "ok" {
		//Now we submit request for leaving chatRoom
		requestLeaveChatRoom := structs.RequestLeaveChatRoom{
			Client:   currentUser,
			ChatRoom: bodyGetChatRoom.ChatRoom,
		}

		//we make the call for getting the chat room
		bLeaveChatRoom := new(bytes.Buffer)
		json.NewEncoder(bLeaveChatRoom).Encode(requestLeaveChatRoom)
		resLeaveChatRoom, _ := http.Post("http://localhost:8888/leaveChatRoom", "application/json; charset=utf-8", bLeaveChatRoom)
		var bodyLeaveChatRoom structs.ResponseLeaveChatRoom
		json.NewDecoder(resLeaveChatRoom.Body).Decode(&bodyLeaveChatRoom)

		if bodyLeaveChatRoom.Status == "ok" {
			fmt.Println("Username " + currentUser.Username + " has left chatRoom " + bodyGetChatRoom.ChatRoom.NameChatRoom)
			fmt.Println()
		}
	} else {
		fmt.Println("There was some error")
	}
}
