package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"

	"../structs"
)

func main() {

	fmt.Println("Starting client...")
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	//Introduce credentials
	fmt.Println("Enter a username: ")
	username := bufio.NewReader(os.Stdin)

	inputUserObject, err := username.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	inputUserObject = inputUserObject[:len(inputUserObject)-1]
	// Asynchronous call

	currentUserNameObject := structs.RequestCreateClient{
		Username: string(inputUserObject)}

	var response structs.ResponseCreateClient

	divCall := client.Go("Clients.CreateClient", currentUserNameObject, &response, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	fmt.Println("Response....")
	if replyCall.Error != nil {
		fmt.Println(replyCall.Error)
	}
	currentGlobalUser := response.Client
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

		switch input {
		case '1':
			createChatRoom(client)
		case '2':
			listChatRoom(client, currentGlobalUser)
		case '3':
			joinChatRoom(client, currentGlobalUser)
		}
	}

}

func createChatRoom(client *rpc.Client) {
	//Input of chatname
	fmt.Print("Choose a name for the chatRoom: ")

	reader := bufio.NewReader(os.Stdin)
	chatName, _ := reader.ReadString('\n')
	chatName = chatName[:len(chatName)-1]

	request := structs.RequestCreateChatRoom{
		NameChatRoom: string(chatName)}

	var response structs.ResponseCreateChatRoom

	divCall := client.Go("ChatRooms.CreateChatRoom", request, &response, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	if replyCall.Error != nil {
		fmt.Println(replyCall.Error)
	}
	fmt.Println(response.Status)

}

//To print in client the chats available
func listChatRoom(client *rpc.Client, currentUser structs.Client) {

	request := structs.RequestListChatRoom{
		Username: string(currentUser.Username)}

	var response structs.ResponseListChatRoom

	divCall := client.Go("ChatRooms.ListChatRoom", request, &response, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	if replyCall.Error != nil {
		fmt.Println(replyCall.Error)
	}
	fmt.Println("Status: " + response.Status)
	fmt.Println("Chats available .....")
	resultArray := response.ChatRooms.Chats
	for _, chatRoom := range resultArray {
		fmt.Print("Name ChatRoom: " + chatRoom.NameChatRoom + " Number of Clients: ")
		fmt.Printf("%d\n", len(chatRoom.Clients.Clients))
		fmt.Println("")
	}

}

//Join some chatRoom
func joinChatRoom(client *rpc.Client, currentUser structs.Client) {

	//Input of chatname
	fmt.Print("Choose a name for joining chatRoom: ")

	reader := bufio.NewReader(os.Stdin)
	chatName, _ := reader.ReadString('\n')
	chatName = chatName[:len(chatName)-1]

	//First we get the chatRoom struct
	requestGetChatRoom := structs.RequestGetChatRoom{
		ChatRoomName: chatName,
	}

	var responseGetChatRoom structs.ResponseGetChatRoom

	divCallGetChatRoom := client.Go("ChatRooms.GetChatRoom", requestGetChatRoom, &responseGetChatRoom, nil)
	replyCallGetChatRoom := <-divCallGetChatRoom.Done // will be equal to divCall
	if replyCallGetChatRoom.Error != nil {
		fmt.Println(replyCallGetChatRoom.Error)
	}

	if responseGetChatRoom.Status == "ok" {
		//Now we submit request for joining chatRoom
		requestJoinChatRoom := structs.RequestJoinChatRoom{
			Client:   currentUser,
			ChatRoom: responseGetChatRoom.ChatRoom,
		}

		var responseJoinChatRoom structs.ResponseJoinChatRoom

		divCallJoinChatRoom := client.Go("ChatRooms.JoinChatRoom", requestJoinChatRoom, &responseJoinChatRoom, nil)
		replyCallJoinChatRoom := <-divCallJoinChatRoom.Done // will be equal to divCall
		if replyCallJoinChatRoom.Error != nil {
			fmt.Println(replyCallJoinChatRoom.Error)
		}
	} else {
		fmt.Println("There was some error")
	}

}
