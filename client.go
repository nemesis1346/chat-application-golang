package main

import(
	"fmt"
	"bufio"
	"os"
	"net"
	"encoding/json"
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
		//we build the object to order the creation of chatroom
		chatRoom:=&ChatRoom{
			userName:"test",
			chatName:"test",
		}
		createChatRoom(chatRoom)

	}

}

func createChatRoom(chatRoom *ChatRoom){
	//TODO PORT MUST BE DYNAMICALLY ADDED
	connection, err:=net.Dial("tcp","localhost:12346")
	if err!=nil{
		fmt.Println(err)
	}
	message,err:=json.Marshal(chatRoom)
	if err!=nil{
		fmt.Println(err)
	}
	connection.Write(message)
}


type ChatRoom struct{
	chatName string
	userName string
}

type Message struct{

}

type ClientModel struct{
	username string
	socket net.Conn
	data chan[] byte
}
