package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"../structs"
)

func main() {
	//Introduce credentials
	fmt.Println("Enter a username: ")

	//optBinBuffer := new(bytes.Buffer)
	conn, err := net.Dial("tcp", "localhost:12346")
	if err != nil {
		fmt.Println(err)
	}

	username := bufio.NewReader(os.Stdin)
	inputUserObject, err := username.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	inputUserObject = inputUserObject[:len(inputUserObject)-1]

	//We create request for create client
	//We send the option over the connection
	mapCreateClient := make(map[string]string)
	mapCreateClient["Username"] = string(inputUserObject)

	optionCreateClient := structs.OptionMessage{
		Option: "1",
		Data:   mapCreateClient,
	}
	gobRequestCreateClient := gob.NewEncoder(conn)
	gobRequestCreateClient.Encode(optionCreateClient)

	flagResponse := false
	//We listen the response
	for {
		response := new(structs.OptionMessage)
		//create a decoder object
		gobResponseOption := gob.NewDecoder(conn)
		error := gobResponseOption.Decode(response)
		if error != nil {
			fmt.Println(error)
		}
		if response.Data["Status"] == "ok" {
			fmt.Println("User " + string(inputUserObject) + " created ...")
			fmt.Println()
			flagResponse = true
			break
		} else {
			fmt.Println("Error: " + response.Data["Status"])
			flagResponse = false
			break
		}
	}
	if flagResponse {
		for {
			//Interface for options of the client
			fmt.Println("Choose an action:")
			fmt.Println("1.Create a chatroom")
			fmt.Println("2.List all existing chatrooms ")
			fmt.Println("3.Join a chatroom ")

			fmt.Print("Choose option: ")

			//Introduce options
			option := bufio.NewReader(os.Stdin)
			//reading the input
			inputOption, _, err := option.ReadRune()
			if err != nil {
				fmt.Println(err)
			}

			switch inputOption {
			case '1':
				createChatRoom()
			case '2':
				listChatRoom()
			case '3':
				joinChatRoom()
			}
		}
		conn.Close()
	}
}

func createChatRoom() {
	//we call the connection
	conn, err := net.Dial("tcp", "localhost:12346")
	if err != nil {
		fmt.Println(err)
	}
	//Input of chatname
	fmt.Print("Choose a name for the chatRoom: ")

	reader := bufio.NewReader(os.Stdin)
	chatName, _ := reader.ReadString('\n')
	chatName = chatName[:len(chatName)-1]

	//We create the map
	mapCreateChatRoom := make(map[string]string)
	mapCreateChatRoom["NameChatRoom"] = string(chatName)

	//We send the option over the connection
	optionMessage := structs.OptionMessage{
		Option: "2",
		Data:   mapCreateChatRoom,
	}
	gobRequestOption := gob.NewEncoder(conn)
	gobRequestOption.Encode(optionMessage)
	//We listen the response
	for {
		response := new(structs.OptionMessage)
		//create a decoder object
		gobResponseOption := gob.NewDecoder(conn)
		error := gobResponseOption.Decode(response)
		if error != nil {
			fmt.Println(error)
		}
		if response.Data["Status"] == "ok" {
			fmt.Println()
			fmt.Println("ChatRoom: " + chatName + " created ...")
			fmt.Println()
			break
		} else {
			fmt.Println("Error: " + response.Data["Status"])
			break
		}
	}
}

func listChatRoom() {
	//we call the connection
	conn, err := net.Dial("tcp", "localhost:12346")
	if err != nil {
		fmt.Println(err)
	}
	//we make the request
	mapListChatRoom := make(map[string]string)

	optionMessage := structs.OptionMessage{
		Option: "3",
		Data:   mapListChatRoom,
	}
	gobReqListChatRoom := gob.NewEncoder(conn)
	gobReqListChatRoom.Encode(optionMessage)
	//we listen to the response
	for {
		response := new(structs.OptionMessage)
		//create a decoder object
		gobResponse := gob.NewDecoder(conn)
		error := gobResponse.Decode(response)
		if error != nil {
			fmt.Println(error)
		}
		fmt.Println("Chats available....")
		fmt.Println()
		if response.Data["Status"] == "ok" {
			delete(response.Data, "Status")
			chatRooms := response.Data
			for k, v := range chatRooms {
				fmt.Println("NameChatRoom:", k, "NumberClients:", v)
				fmt.Println()
			}
			break
		} else {
			fmt.Println("Error: " + response.Data["Status"])
			break
		}
	}
}
func joinChatRoom() {
	//First we show the chats available
	listChatRoom()
	//we call the connection
	conn, err := net.Dial("tcp", "localhost:12346")
	if err != nil {
		fmt.Println(err)
	}

	//Input of chatname
	fmt.Print("Choose a name for joining chatRoom: ")

	reader := bufio.NewReader(os.Stdin)
	chatName, _ := reader.ReadString('\n')
	chatName = chatName[:len(chatName)-1]

	//we make the request
	mapJoinChatRoom := make(map[string]string)

	mapJoinChatRoom["NameChatRoom"] = string(chatName)

	optionMessage := structs.OptionMessage{
		Option: "4",
		Data:   mapJoinChatRoom,
	}
	gobReqJoinChatRoom := gob.NewEncoder(conn)
	gobReqJoinChatRoom.Encode(optionMessage)
	//we listen the response
	for {
		response := new(structs.OptionMessage)
		//create a decoder object
		gobResponse := gob.NewDecoder(conn)
		error := gobResponse.Decode(response)
		if error != nil {
			fmt.Println(error)
		}
		if response.Data["Status"] == "ok" {
			fmt.Println("Client joined to  " + string(chatName))
			break
		} else {
			fmt.Println("Error: " + response.Data["Status"])
			break
		}
	}
}

//Get messages
func getMessages() {

}
