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
			createChatRoom(conn)
		case '2':
			listChatRoom(conn)
		case '3':
			joinChatRoom(conn)
		}
	}
	conn.Close()

}

func createChatRoom(conn net.Conn) {
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
		gobRequestOption := gob.NewDecoder(conn)
		error := gobRequestOption.Decode(response)
		if error != nil {
			fmt.Println(error)
		}
		if response.Data["Status"] == "ok" {
			fmt.Println("ChatRoom: " + chatName + " created ...")
			break
		} else {
			fmt.Println("Error: " + response.Data["Status"])
			break
		}
	}
}

func listChatRoom(conn net.Conn) {

}
func joinChatRoom(conn net.Conn) {

}
