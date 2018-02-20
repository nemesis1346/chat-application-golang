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

	divCall := client.Go("ChatRooms.CreateClient", request, &response, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	fmt.Println(replyCall.Reply)
}

func listChatRoom(client *rpc.Client, currentUser structs.Client) {
	request := structs.RequestListChatRoom{
		Username: string(currentUser.Username)}

	var response structs.ResponseCreateChatRoom

	divCall := client.Go("ChatRooms.ListChatRoom", request, &response, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	fmt.Println(replyCall.Reply)
}
