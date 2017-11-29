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
	inputUserObject=inputUserObject[:len(inputUserObject)-1]
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
	//Input for creating a new chatroom
	fmt.Print("Choose a name for the chatRoom: ")
	reader:=bufio.NewReader(os.Stdin)
	chatName,_:=reader.ReadString('\n')
	chatName=chatName[:len(chatName)-1]
	//We create the object
	//chatOrder:=ChatStruct{ ChatName:chatName}
	//chatOrderJson,err:=json.Marshal(chatOrder)
	//if err!=nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(string(chatOrderJson))

	//Create message general
	jsonContent:= OptionMessageClient{
		"1",
		string(userObject.Username),
		string(chatName)}

	message,err:=json.Marshal(jsonContent)
	if err!=nil{
		fmt.Println(err)
	}

	conn.Write([]byte(strings.TrimRight(string(message),"\n")))
}

func listChatRoom(conn net.Conn, userObject UsernameStruct){
	//Create the instruction message for list all the chatrooms
	jsonContent:=OptionMessageClient{"2",
		string(userObject.Username),""}

	message,err:=json.Marshal(jsonContent)
	if err!=nil{
		fmt.Println(err)
	}
	conn.Write([]byte(strings.TrimRight(string(message),"\n")))
}

func joinChatRoom(conn net.Conn, userObject UsernameStruct){
	//Indicate to list all the chatrooms first to join
	jsonContent:=OptionMessageClient{"3",
		string(userObject.Username),""}
	message,err:=json.Marshal(jsonContent)
	if err!=nil{
		fmt.Println(err)
	}
	conn.Write([]byte(strings.TrimRight(string(message),"\n")))

	//Input to choose the chatroom to join
	fmt.Print("Introduce the name of the chatRoom you want to join: ")
	reader:=bufio.NewReader(os.Stdin)
	chatName,_:=reader.ReadString('\n')
	chatName=chatName[:len(chatName)-1]
	//Create the instruction message for joining a specific chatroom
	jsonContent2:=OptionMessageClient{"3",
		string(userObject.Username),chatName}
	message2,err:=json.Marshal(jsonContent2)
	if err!=nil{
		fmt.Println(err)
	}
	conn.Write([]byte(strings.TrimRight(string(message2),"\n")))

	//Start to write messages
	fmt.Println("Start sending messages......")
	client:=&Client{socket:conn}
	go client.receive()
	for{
		reader:=bufio.NewReader(os.Stdin)
		message,_:=reader.ReadString('\n')
		message=message[:len(message)]
		//TODO here i should create a username of the client fmt.Println(username);

		conn.Write([]byte(strings.TrimRight(string(message),"\n")))
	}

}
func(client *Client) receive(){
	for{
		message:=make([]byte, 4096)
		length, err:=client.socket.Read(message)
		if err!=nil{
			client.socket.Close()
			break
		}
		if length>0{
			fmt.Println("RECEIVED: "+string(message))
		}
	}
}

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

type Client struct{
	socket net.Conn
	data chan[]byte
}
