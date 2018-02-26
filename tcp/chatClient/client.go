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
	username := bufio.NewReader(os.Stdin)

	inputUserObject, err := username.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	inputUserObject = inputUserObject[:len(inputUserObject)-1]
	// userObject := UsernameStruct{
	// 	Username: string(inputUserObject)}

	// create a temp buffer
	// tmp := make([]byte, 500)

	// for {

	// 	_, err = connOption.Read(tmp)
	// 	fmt.Println(tmp)
	// 	tmpbuff := bytes.NewBuffer(tmp)
	// 	fmt.Println(tmpbuff)

	// 	tmpstruct := new(structs.OptionMessage)

	// 	// creates a decoder object
	// 	gobobj := gob.NewDecoder(connOption)

	// 	// decodes buffer and unmarshals it into a Message struct
	// 	gobobj.Decode(tmpstruct)

	// 	// lets print out!
	// 	fmt.Println(tmpstruct)
	// }
	//bytesBeingSent := optBinBuffer.Bytes()
	//bytesBeingSent = bytesBeingSent[:len(bytesBeingSent)-1]
	//fmt.Println(bytesBeingSent)

	//bytesBeingSent = bytes.Trim(bytesBeingSent, "\x00")
	//connOption.Write(bytesBeingSent)

	for {
		//Interface for options of the client
		fmt.Println("Choose an action:")
		fmt.Println("1.Create a chatroom")
		fmt.Println("2.List all existing chatrooms ")
		fmt.Println("3.Join a chatroom ")
		fmt.Println("4.Leave a chatroom \n")

		fmt.Print("Choose option: ")

		//Introduce options
		option := bufio.NewReader(os.Stdin)
		//reading the input
		inputOption, _, err := option.ReadRune()
		if err != nil {
			fmt.Println(err)
		}

		//We send the option over the connection
		optionMessage := structs.OptionMessage{
			Option: string(inputOption),
		}
		//optBinBuffer := new(bytes.Buffer)
		conn, err := net.Dial("tcp", "localhost:12346")
		if err != nil {
			fmt.Println(err)
		}
		gobRequestOption := gob.NewEncoder(conn)
		gobRequestOption.Encode(optionMessage)

		// 	//TODO PORT MUST BE DYNAMICALLY ADDED
		// 	connection, err := net.Dial("tcp", "localhost:12346")
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}

		switch inputOption {
		case '1':
			createChatRoom(conn)
		case '2':
			listChatRoom(conn)
		case '3':
			joinChatRoom(conn)
		}
		conn.Close()
	}

}

func createChatRoom(conn net.Conn) {

}

func listChatRoom(conn net.Conn) {

}
func joinChatRoom(conn net.Conn) {

}
