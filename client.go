package main

import(
	"fmt"
	"bufio"
	"os"
	"net"
)


func main(){
	reader:=bufio.NewReader(os.Stdin)

	//Interface for options of the client
	fmt.Println("Choose an action:")
	fmt.Println("1.Create a chatroom")
	fmt.Println("2.List all existing chatrooms ")
	fmt.Println("3.Join a chatroom ")
	fmt.Println("3.Send messages to chat rooms")
	fmt.Println("5.Leave a chatroom ")

	//reading the input
	input,_,err:=reader.ReadRune()
	if err!=nil{
		fmt.Println(err)
	}

	switch input {
	case '1':

	}

}

func createChatroom(){
	//TODO PORT MUST BE DYNAMICALLY ADDED
	connection, err:=net.Dial("tcp","localhost:12345")
	if err!=nil{
		fmt.Println(err)
	}

}

type ChatRoom struct{

}

type Message struct{

}

type ClientModel struct{
	username string
	socket net.Conn
	data chan[] byte
}
