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

			var msg ChatRoom

			//We decode the byte array that contains the json
			err=json.Unmarshal([]byte(string(jsonResult)), &msg)
			if err!=nil{
				fmt.Println(err)
			}

			fmt.Println("jsonResponse: "+string(jsonResult))
			fmt.Println(msg.ChatName)
			conn.Close()
		}
		conn.Close()
	}
}
type ChatRoom struct{
	ChatName string `json:"chatName"`
	UserName string `json:"userName"`
}
