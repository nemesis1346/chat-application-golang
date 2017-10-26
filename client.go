package main

import(
	"fmt"
	"bufio"
	"os"
	"net"
	"encoding/json"
	"strings"
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
		//chatRoom:=&ChatRoom{
		//	userName:"test",
		//	chatName:"test",
		//}
		createChatRoom()

	}

}

func createChatRoom(){
	//TODO PORT MUST BE DYNAMICALLY ADDED
	connection, err:=net.Dial("tcp","localhost:12346")
	if err!=nil{
		fmt.Println(err)
	}
	//new try
	contentToSend:=Message{"test","test"}
	message,err:=json.Marshal(contentToSend)
	if err!=nil{
		fmt.Println(err)
	}
	//fmt.Println(message)
	//fmt.Println("Message send: "+string(message))
	//fmt.Println(len([]byte(message)))
	connection.Write([]byte(strings.TrimRight(string(message),"\n")))
}


type Message struct{
	ChatName string `json:"chatName"`
	UserName string `json:"userName"`
}

type ClientModel struct{
	username string
	socket net.Conn
	data chan[] byte
}
