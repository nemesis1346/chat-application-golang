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

	requestNewUser := structs.RequestCreateClient{
		Username: string(inputUserObject)}

	var response structs.ResponseCreateClient

	divCall := client.Go("Clients.CreateClient", requestNewUser, &response, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	if replyCall.Error != nil {
		fmt.Println(replyCall.Error)
	}
	if response.Status == "ok" {
		fmt.Println()
		fmt.Println("----User: " + response.Client.Username + " created!!")
		fmt.Println()
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
			case '4':
				leaveChatRoom(client, currentGlobalUser)
			}
		}
	} else {
		fmt.Println()
		fmt.Println(response.Status)
		fmt.Println()
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
	if response.Status == "ok" {
		fmt.Println()
		fmt.Println("----ChatCreated: " + response.Status)
		fmt.Println("NameChatRoom: " + response.ChatRoom.NameChatRoom)
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println("----Error: " + response.Status)
		fmt.Println()
	}

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
	fmt.Println()
	fmt.Println("Chats available .....")
	fmt.Println()

	resultArray := response.ChatRooms.Chats

	if len(resultArray) > 0 {
		for _, chatRoom := range resultArray {
			fmt.Print("Name ChatRoom: " + chatRoom.NameChatRoom + ", Number of Clients: ")
			fmt.Printf("%d\n", len(chatRoom.Clients.Clients))
			fmt.Println("")
		}
		fmt.Println()
	} else {
		fmt.Println("There is no chats available")
		fmt.Println()
	}

}

//Join some chatRoom
func joinChatRoom(client *rpc.Client, currentUser structs.Client) {
	//First we show the chats available
	fmt.Print("List of Chats available.... ")
	listChatRoom(client, currentUser)

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
		fmt.Println()
		fmt.Println("Client " + currentUser.Username + " joined to " + requestJoinChatRoom.ChatRoom.NameChatRoom)
		fmt.Println()

		//Now we get the rest of the messages
		fmt.Println("Previous messages....")
		requestGetPreviousMessages := structs.RequestGetPreviousMessages{
			ChatRoom: responseGetChatRoom.ChatRoom,
			Client:   currentUser,
		}
		//var
		//Now we start chating
		fmt.Println("Start chating.....")
		if responseJoinChatRoom.Status == "ok" {
			for {

				readerMessage := bufio.NewReader(os.Stdin)
				messageContent, _ := readerMessage.ReadString('\n')
				messageContent = chatName[:len(chatName)-1]
				//We submit request for the message
				requestSaveMessage := structs.RequestSaveMessage{
					Client:   currentUser,
					Content:  messageContent,
					ChatRoom: responseGetChatRoom.ChatRoom,
				}
				var responseSaveMessage structs.ResponseSaveMessage

				divCallSaveMessage := client.Go("ChatRooms.JoinChatRoom", requestSaveMessage, &responseSaveMessage, nil)
				replyCallSaveMesssage := <-divCallSaveMessage.Done // will be equal to divCall
				if replyCallSaveMesssage.Error != nil {
					fmt.Println(replyCallSaveMesssage.Error)
				}
				if responseSaveMessage.Status == "ok" {
					fmt.Println("Username: " + responseSaveMessage.Content + " delivered")
				} else {
					fmt.Println("There was an error in saving message")
				}
			}
		}

	} else {
		fmt.Println("There was some error")
	}

}

//This method must be called when I am messaging
func leaveChatRoom(client *rpc.Client, currentUser structs.Client) {

	//First we show the user the active ChatRooms
	fmt.Println("List of ChatRooms active")
	listChatRoom(client, currentUser)

	//Input of chatname to leave
	fmt.Print("Choose a name for leaving chatRoom: ")

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
		//Now we submit request for leaving chatRoom
		requestLeaveChatRoom := structs.RequestLeaveChatRoom{
			Client:   currentUser,
			ChatRoom: responseGetChatRoom.ChatRoom,
		}
		var responseLeaveChatRoom structs.ResponseLeaveChatRoom

		divCallLeaveChatRoom := client.Go("ChatRooms.LeaveChatRoom", requestLeaveChatRoom, &responseLeaveChatRoom, nil)
		replyCallLeaveChatRoom := <-divCallLeaveChatRoom.Done // will be equal to divCall
		if replyCallLeaveChatRoom.Error != nil {
			fmt.Println(replyCallLeaveChatRoom.Error)
		}

	} else {
		fmt.Println("There was some error")
	}
}

func sendMessage(client *rpc.Client, currentUser structs.Client,
	currentChat structs.ChatRoom, message string) {

}
