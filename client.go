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
	//Introduce credentials
	fmt.Println("Enter a username: ")
	username:=bufio.NewReader(os.Stdin)

	inputUserObject,err:=username.ReadString('\n')
	if err!=nil{
		fmt.Println(err)
	}
	userObject:=UsernameStruct{
		Username:string(inputUserObject)}

	for{

		//Introduce options
		option:=bufio.NewReader(os.Stdin)

		//Interface for options of the client
		fmt.Println("Choose an action:")
		fmt.Println("1.Create a chatroom")
		fmt.Println("2.List all existing chatrooms ")
		fmt.Println("3.Join a chatroom ")
		fmt.Println("4.Leave a chatroom \n")

		fmt.Print("Choose option: ")
		//reading the input
		input,_,err:=option.ReadRune()
		if err!=nil{
			fmt.Println(err)
		}

		//TODO PORT MUST BE DYNAMICALLY ADDED
		connection, err:=net.Dial("tcp","localhost:12346")
		if err!=nil{
			fmt.Println(err)
		}

		switch input {
		case '1':
			createChatRoom(connection, userObject)
		case '2':
			listChatRoom(connection,userObject)
		case '3':
			joinChatRoom(connection,userObject)
		case '4':
			leaveChatRoom(connection,userObject)
		}
	}

}

func createChatRoom(conn net.Conn, userObject UsernameStruct){

	chatOrder:=ChatStruct{ ChatName:"test"}
	chatOrderJson,err:=json.Marshal(chatOrder)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(string(chatOrderJson))

	//Create message general
	jsonContent:= OptionMessageClient{
		"1",
		string(userObject.Username),
		string(chatOrderJson)}

	message,err:=json.Marshal(jsonContent)
	if err!=nil{
		fmt.Println(err)
	}
	//fmt.Println(message)
	//fmt.Println("OptionMessageClient send: "+string(message))
	//fmt.Println(len([]byte(message)))
	conn.Write([]byte(strings.TrimRight(string(message),"\n")))
}
func listChatRoom(conn net.Conn, userObject UsernameStruct){

}

func joinChatRoom(conn net.Conn, userObject UsernameStruct){}

func leaveChatRoom(conn net.Conn, userObject UsernameStruct){}

type UsernameStruct struct{
	Username string `json:"username"`
}

type OptionMessageClient struct{
	Option string `json:"option"`
	UserName string `json:"userName"`
	Data string `json:"data"`
}

type ChatStruct struct{
	ChatName string `json:"chatName"`
}
