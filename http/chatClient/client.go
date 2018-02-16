package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
	//HTTP POST OR GET
	resp, err := http.Get("http://localhost:8888/endpoint1")
	if err != nil {
		//handle error
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(keepLines(string(body), 3))
}
func listChatRoom(userObject UsernameStruct) {

}
func joinChatRoom(userObject UsernameStruct) {

}
func leaveChatRoom(userObject UsernameStruct) {

}

func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

type UsernameStruct struct {
	Username string `json:"username"`
}
