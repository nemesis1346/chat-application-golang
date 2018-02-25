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

	requestOption := new(structs.OptionMessage)

	//create a decoder object
	gobRequestOption := gob.NewDecoder(conn)

	error := gobRequestOption.Decode(requestOption)
	fmt.Println(error)
	fmt.Println(requestOption.Option)

	conn.Close()

}

// func createClient() {
// 	//request
// 	var requestCreateClient structs.requestCreateClient
// 	gobRequestCreateClient := gob.NewDecoder(requestCreateClient)
// }
