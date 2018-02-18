package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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
	fmt.Println("----CREATE CHAT ROOM")
	var c ClientObject
	req, err := c.newRequest("POST", "/createChatRoom", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(req)

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
	fmt.Println("----LIST CHAT ROOM")
	var c ClientObject
	req, err := c.newRequest("GET", "/listChatRoom", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(req)
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

//This is a generic method to generate a request
func (c *ClientObject) newRequest(method string, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{
		Host: "http://localhost:8888",
		Path: path,
	}
	fmt.Println(c)

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	//Figure it out why do I need User Agent?
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

//This is a generic method to make the calls
func (c *ClientObject) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
