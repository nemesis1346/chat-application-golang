package main

import (
	"encoding/gob"
	"fmt"
	"net"

	"../structs"
)

func main() {

	//We create a listener
	fmt.Println("Starting server...")
	fmt.Println("-----------------------")
	listener, err := net.Listen("tcp", ":12346")
	if err != nil {
		fmt.Println(err)
	}
	//We accept the connection
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		conn.Close()
	}

	//now we create new client

	//We create a buffer
	// tmp := make([]byte, 1024)
	// fmt.Println(tmp)

	// _, err = conn.Read(tmp)
	// fmt.Println(tmp)

	//We take off extra zeros
	//tmp = bytes.Trim(tmp, "\x00")

	//fmt.Println(tmp)

	//Temporal buffer for incomming messages
	//tmpBuffer := bytes.NewBuffer(tmp)

	//Incomming option messages
	for {
		requestOption := new(structs.OptionMessage)
		//create a decoder object
		gobRequestOption := gob.NewDecoder(conn)
		error := gobRequestOption.Decode(requestOption)
		fmt.Println(error)
		fmt.Println(requestOption.Option)
		switch requestOption.Option {
		case "1":
			createChatRoom(conn)
		case "2":
			listChatRoom(conn)
		case "3":
			joinChatRoom(conn)
		}
		conn.Close()
	}
	conn.Close()
}

func createChatRoom(conn net.Conn) {

	fmt.Print("create chat room")
	conn.Close()

}
func listChatRoom(conn net.Conn) {

}
func joinChatRoom(conn net.Conn) {

}

// func createClient() {
// 	//request
// 	var requestCreateClient structs.requestCreateClient
// 	gobRequestCreateClient := gob.NewDecoder(requestCreateClient)
// }
