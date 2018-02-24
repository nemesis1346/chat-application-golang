package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"../structs"
)

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
	fmt.Println(resCreateClient)
	json.NewDecoder(resCreateClient.Body).Decode(&bodyCreateClient)
	fmt.Println(bodyCreateClient.Status)

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

}
func joinChatRoom(currentUser structs.Client) {

}
func leaveChatRoom(currentUser structs.Client, chatName string) {

}
