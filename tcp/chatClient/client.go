package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"../structs"
)

//var time
var currentTime time.Time

func main() {
	//Introduce credentials
	fmt.Println("Enter a username: ")

	//optBinBuffer := new(bytes.Buffer)
	conn, err := net.Dial("tcp", "localhost:12346")
	if err != nil {
		fmt.Printf("[Dial]\t", err)
	}

	username := bufio.NewReader(os.Stdin)
	inputUserObject, err := username.ReadString('\n')
	if err != nil {
		fmt.Printf("[RequestOption]\t", err)
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
			fmt.Println("4.Exit Program ")

			fmt.Print("Choose option: ")

			//Introduce options
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("[InputReader]\t", err)
			}

			inputNumber, err := strconv.Atoi(input[0:(len(input) - 1)])
			if err != nil {
				fmt.Printf("[InputNumber]\t", err)
			}

			switch inputNumber {
			case 1:
				createChatRoom()
			case 2:
				listChatRoom()
			case 3:
				joinChatRoom(string(inputUserObject))
			case 4:
				conn.Close()
				os.Exit(0)
			default: // -- TODO: What if the input is not recognized? Provide a default case.
				fmt.Printf("Unknown option: %d\n", inputNumber)
			}

		}
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
		fmt.Printf("[ListenChatRoom Connection]\t", err)
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
func joinChatRoom(username string) {
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
	flagResJoinChatRoom := false
	for {
		response := new(structs.OptionMessage)
		//create a decoder object
		gobResponse := gob.NewDecoder(conn)
		error := gobResponse.Decode(response)
		if error != nil {
			fmt.Println(error)
		}
		if response.Data["Status"] == "ok" {
			fmt.Println()
			fmt.Println("Client joined to  " + string(chatName))
			fmt.Println()
			flagResJoinChatRoom = true
			break
		} else {
			fmt.Println("Error: " + response.Data["Status"])
			break
		}
	}
	if flagResJoinChatRoom {
		getPreviousMessages(string(chatName))

		//current message time TODO: the time should be the last time of the message
		currentTime = time.Now()
		//Now we start chating
		fmt.Println("Start chating.....")
		fmt.Println()

		go listenMessages(username, string(chatName)) //Start listening messages

		//We stay in a loop for chating
		for {
			connChating, err := net.Dial("tcp", "localhost:12346")
			if err != nil {
				fmt.Println(err)
			}

			readerMessage := bufio.NewReader(os.Stdin)
			messageContent, _ := readerMessage.ReadString('\n')
			messageContent = messageContent[:len(messageContent)]
			//First we detect the command exit
			if strings.TrimRight(messageContent, "\n") == "exit" {
				leaveChatRoom(username, string(chatName))
				break
			} else {
				//We submit request for the message
				mapSaveMessage := make(map[string]string)
				mapSaveMessage["Username"] = username
				mapSaveMessage["Content"] = messageContent
				mapSaveMessage["NameChatRoom"] = string(chatName)

				requestSaveMessage := structs.OptionMessage{
					Option: "7",
					Data:   mapSaveMessage,
				}

				//we make the call for saving messages
				gobReqSaveMessage := gob.NewEncoder(connChating)
				gobReqSaveMessage.Encode(requestSaveMessage)

				for {
					response := new(structs.OptionMessage)
					//create a decoder object
					gobResponse := gob.NewDecoder(connChating)
					err := gobResponse.Decode(response)
					if err != nil {
						fmt.Printf("[DecodeResponseGob-ListenMessageResponse]\t", err)
					}
					if response.Data["Status"] != "ok" {
						fmt.Println("There was an error in saving message")
						break
					} else {
						//We update the time stamp of the last message
						layoutTimeLast := time.RFC3339
						currentTime, err = time.Parse(layoutTimeLast, response.Data["TimeLast"])
						if err != nil {
							fmt.Printf("[ListenMessageResponse]\t", err)
						}
						break
					}
				}
			}
		}
	}

}

//Listen constantly messages
func listenMessages(username string, nameChatRoom string) {
	for {
		//we call the connection
		conn, err := net.Dial("tcp", "localhost:12346")
		if err != nil {
			fmt.Printf("[ListenMessages]\t", err)
		}
		//we make the request
		mapListenMessages := make(map[string]string)
		mapListenMessages["Time"] = currentTime.Format(time.RFC3339)
		mapListenMessages["Username"] = username
		mapListenMessages["NameChatRoom"] = nameChatRoom
		optionMessage := structs.OptionMessage{
			Option: "8",
			Data:   mapListenMessages,
		}
		gobReqListenMessages := gob.NewEncoder(conn)
		gobReqListenMessages.Encode(optionMessage)
		//We get the previous messages
		for {
			response := new(structs.OptionMessage)
			//create a decoder object
			gobResponse := gob.NewDecoder(conn)
			error := gobResponse.Decode(response)
			if error != nil {
				fmt.Println(error)
			}
			fmt.Println(len(response.Data))
			if response.Data["Status"] == "ok" {
				delete(response.Data, "Status")
				messages := response.Data
				for k, v := range messages {
					fmt.Println(k + ": " + " " + v)
				}
				currentTime = time.Now()
			}
			break
		}
		time.Sleep(time.Second / 2)
	}
}

//Leave chat room
func leaveChatRoom(username string, nameChatRoom string) {
	//we call the connection
	conn, err := net.Dial("tcp", "localhost:12346")
	if err != nil {
		fmt.Println(err)
	}

	//we make the request
	mapLeaveChatRoom := make(map[string]string)

	mapLeaveChatRoom["NameChatRoom"] = nameChatRoom
	mapLeaveChatRoom["Username"] = username

	optionMessage := structs.OptionMessage{
		Option: "6",
		Data:   mapLeaveChatRoom,
	}
	gobReqLeaveMessage := gob.NewEncoder(conn)
	gobReqLeaveMessage.Encode(optionMessage)
	//We get the previous messages
	for {
		response := new(structs.OptionMessage)
		//create a decoder object
		gobResponse := gob.NewDecoder(conn)
		error := gobResponse.Decode(response)
		if error != nil {
			fmt.Println(error)
		}
		if response.Data["Status"] == "ok" {
			fmt.Println("Username " + username + " left the chat room")
		} else {
			fmt.Println(response.Data["Status"])
		}
	}
}

//Get previous messages
func getPreviousMessages(nameChatRoom string) {
	//we call the connection
	conn, err := net.Dial("tcp", "localhost:12346")
	if err != nil {
		fmt.Println(err)
	}

	//we make the request
	mapGetPreviousMessages := make(map[string]string)

	mapGetPreviousMessages["NameChatRoom"] = nameChatRoom

	optionMessage := structs.OptionMessage{
		Option: "5",
		Data:   mapGetPreviousMessages,
	}
	gobReqGetPreviousMessages := gob.NewEncoder(conn)
	gobReqGetPreviousMessages.Encode(optionMessage)
	//We get the previous messages
	for {
		response := new(structs.OptionMessage)
		//create a decoder object
		gobResponse := gob.NewDecoder(conn)
		error := gobResponse.Decode(response)
		if error != nil {
			fmt.Println(error)
		}
		if response.Data["Status"] == "ok" {
			delete(response.Data, "Status")
			message := response.Data
			for k, v := range message {
				fmt.Println("Username:", k, " message: ", v)
				fmt.Println()
			}
			break
		} else {
			fmt.Println(response.Data["Status"])
			break
		}
	}
}
