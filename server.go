package main
import (
	"fmt"
	"net"
	"encoding/json"
	"bytes"
)
func main(){
	startServer()
}
func startServer(){

	//Here we initialize the slice of chatrooms(ServerManagers) , initializing of spaces in slice in 1
	chatrooms:=make([]*ChatRoomManagerServer,1)

	fmt.Println("Starting server...")
	//We create a listener
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
			fmt.Println(msg.UserName)
			fmt.Println(msg.Data)

			//MANAGE THE OPTIONS
			switch msg.Option {
			case "1":
				createChatRoomServer(conn,msg.Data,chatrooms)
			case "2":
				listChatRoomServer(conn)
			case "3":
				joinChatRoomServer(conn)
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
func createChatRoomServer(conn net.Conn, data string,chatrooms []*ChatRoomManagerServer){

	//We create chatRoomStruct struct
	//data is the name of the server in this option
	//currentChatRoom:=ChatRoomManagerServer{nameChatRoom:data}
	//append(chatrooms, currentChatRoom)

}
func listChatRoomServer(conn net.Conn){}
func joinChatRoomServer(conn net.Conn){}
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