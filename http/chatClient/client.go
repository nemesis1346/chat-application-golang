package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	userObject := UsernameStruct{
		Username: string(inputUserObject)}

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
			createChatRoom(userObject)
		case '2':
			listChatRoom(userObject)
		case '3':
			joinChatRoom(userObject)
		case '4':
			leaveChatRoom(userObject)
		}
	}

}

func createChatRoom(userObject UsernameStruct) {
	//Input of chatname
	fmt.Print("Choose a name for the chatRoom: ")

	reader := bufio.NewReader(os.Stdin)
	chatName, _ := reader.ReadString('\n')
	chatName = chatName[:len(chatName)-1]

	request := structs.RequestCreateChatRoom{
		NameChatRoom: string(chatName)}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(request)

	res, _ := http.Post("http://localhost:8888/createChatRoom", "application/json; charset=utf-8", b)

	var body structs.ResponseCreateChatRoom
	json.NewDecoder(res.Body).Decode(&body)
	fmt.Println(body.ChatRoom.NameChatRoom)
}

func listChatRoom(userObject UsernameStruct) {
	//HTTP POST OR GET
	// resp, err := http.Get("http://localhost:8888/endpoint1")
	// if err != nil {
	// 	handle error
	// }
	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(keepLines(string(body), 3))
	// fmt.Println("----LIST CHAT ROOM")
	// var c ClientObject
	// req, err := c.newRequest("GET", "/listChatRoom", nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(req)
}
func joinChatRoom(userObject UsernameStruct) {

}
func leaveChatRoom(userObject UsernameStruct) {

}

type UsernameStruct struct {
	Username string `json:"username"`
}

type ClientObject struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}
