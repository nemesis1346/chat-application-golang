package main
import (
	"fmt"
	"net"
	"encoding/json"
	"bytes"
	"strings"
)
func main(){
	startServer()
}
func startServer(){
	//Here we initialize the slice of chatrooms(ServerManagers) , initializing of spaces in slice in 1
	var chatrooms []*ChatRoomManagerServer

	//We create a listener
	fmt.Println("Starting server...")
	fmt.Println("-----------------------")

	listener, err:=net.Listen("tcp",":12346")
	if err!=nil{
		fmt.Println(err)
	}
	for{
		//We accept the connection
		conn,err:=listener.Accept()
		if err!=nil{
			fmt.Println(err)
			conn.Close()
		}
		//We create a buffer
		buffer:=make([]byte,1024)

		length,err:=conn.Read(buffer)
		if err!=nil{
			fmt.Println(err)
			conn.Close()
		}
		if length>0{
			//manage json
			//We are filtering the nul points of the message to a json byte array
			jsonResult:=bytes.Trim(buffer,"\x00")
			//fmt.Println(correctMesage)

			var msg OptionMessageServer

			//We decode the byte array that contains the json
			err=json.Unmarshal([]byte(string(jsonResult)), &msg)
			if err!=nil{
				fmt.Println(err)
			}

			fmt.Println("jsonResponse: "+string(jsonResult))
			fmt.Println("-----------------------")

			//MANAGE THE OPTIONS
			switch msg.Option {
			case "1":
				temporaryChatRoom:=createChatRoomServer(conn,msg.Data,chatrooms)
				chatrooms=append(chatrooms, temporaryChatRoom)
				fmt.Println("New chatRoom Created: ",temporaryChatRoom.nameChatRoom)
				fmt.Println("-----------------------")

			case "2":
				listChatRoomServer(conn, chatrooms)
			case "3":
				joinChatRoomServer(conn,chatrooms)
			case "4":
				leaveChatRoomServer(conn)
			}
			conn.Close()
		}
		conn.Close()
	}
}
type ChatRoom struct{
	ChatName string `json:"chatName"`
	UserName string `json:"userName"`
}
func createChatRoomServer(conn net.Conn, data string,chatrooms []*ChatRoomManagerServer)*ChatRoomManagerServer{
	//We create chatRoomStruct struct
	//data is the name of the server in this option
	currentChatRoom:=&ChatRoomManagerServer{nameChatRoom:data}
	return currentChatRoom
}

func printSlice(s []*ChatRoomManagerServer) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}


func listChatRoomServer(conn net.Conn, chatrooms[]*ChatRoomManagerServer){
	//printSlice(chatrooms)
	fmt.Println("LIST OF CHATS")

	for _, chatroom:=range chatrooms{
		fmt.Println("-------------------------")
		if chatroom.nameChatRoom!=""{
			fmt.Println("Name: ", chatroom.nameChatRoom)
		}
		if chatroom.clients!=nil{
			fmt.Println("Clients: ",chatroom.clients)
		}
		fmt.Println("-------------------------")
	}
}
func joinChatRoomServer(conn net.Conn,chatrooms[]*ChatRoomManagerServer){
	//Introduce the name of the chat you want to join
	listChatRoomServer(conn, chatrooms)
	fmt.Println("SELECT THE NAME OF THE CHAT YOU WANT TO JOIN")


		//We create a buffer
		message:=make([]byte,1024)

		length,err:=conn.Read(message)
		if err!=nil{
			fmt.Println(err)
			conn.Close()
		}
		if length>0 {
			jsonResult:=bytes.Trim(message,"\x00")

			var msg OptionMessageServer

			//We decode the byte array that contains the json
			err=json.Unmarshal([]byte(string(jsonResult)), &msg)
			if err!=nil{
				fmt.Println(err)
			}

			fmt.Println("jsonResponse: "+string(jsonResult))
			fmt.Println("-----------------------")

			for _, chatroom := range chatrooms {
				fmt.Println("msgData: "+msg.Data)
				fmt.Println("list Name: "+chatroom.nameChatRoom)
				if (strings.TrimRight(chatroom.nameChatRoom, "\n") == string(msg.Data)) {
					fmt.Print("CHAT :", chatroom.nameChatRoom)
					fmt.Print(" INITIATED")
					fmt.Println("-------------------------")
					//STARTING CHAT
					//go chatroom.start()
					//for{
					//	message:=make([]byte,1024)
					//	jsonResultMessage:=bytes.Trim(message,"\x00")
					//	fmt.Println(jsonResultMessage)
						//Here is adding the information of the connection to an instance of a client
						//client:=&Client{socket: connection, data: make(chan []byte)}
						//chatroom.register<-client
						//go chatroom.receive(client)
						//go chatroom.send(client)
					//}
				} else {
					fmt.Println("CHAT DOESNT EXIST")
				}
			}

	}

}
func leaveChatRoomServer(conn net.Conn){}


type OptionMessageServer struct{
	Option string `json:"option"`
	UserName string `json:"userName"`
	Data string `json:"data"`
}

type ClientServer struct{
	socket net.Conn
	infoClient chan[]byte
	messages chan[]byte
}

type ChatRoomManagerServer struct{
	nameChatRoom string
	clients map[*ClientServer]bool
	broadcast chan[]byte
	register chan *ClientServer
	unregister chan *ClientServer
}

func(manager *ChatRoomManagerServer) start(){
	for{
		select{
		case connection:=<-manager.register:
			manager.clients[connection]=true
			fmt.Println("Added new connection")
		case connection:=<-manager.unregister:
			if _,ok:=manager.clients[connection];ok{
				//TODO maybe here it is erasing the data, take a look to it
				close(connection.messages)
				close(connection.infoClient)
				fmt.Println("A connection has terminated")
			}
		case message:=<-manager.broadcast:
			for connection:=range manager.clients{
				select{
				case connection.messages<-message:
				default:
					//TODO maybe here it is erasing the data, take a look to it
					close(connection.messages)
					close(connection.infoClient)
					delete(manager.clients, connection)
				}
			}
		}
	}
}