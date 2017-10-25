package main

import (
	"fmt"
	"net"

	"encoding/json"
)

func main(){
	startServer()
}


func startServer(){
	fmt.Println("Starting server...")
	listener, err:=net.Listen("tcp",":12346")
	if err!=nil{
		fmt.Println(err)
	}
	for{
		conn,err:=listener.Accept()
		if err!=nil{
			fmt.Println(err)
			conn.Close()
		}
		messageResult:=make([]byte,4096)
		var messageResult2 map[string]interface{}
		length,err:=conn.Read(messageResult)
		if length>0{
			//manage json
			fmt.Print("respuesta:"+string(messageResult))
		}

		//prueba:=ChatRoom{
		//	chatName: "",
		//	userName: "",
		//}
		//json.Unmarshal(messageResult, messageResult2)
		json.Unmarshal([]byte(string(messageResult)),&messageResult2)
		//decode(conn)

		fmt.Println("%v", messageResult2)
		conn.Close()

	}
}

//func decode(message io.Reader)(result *ChatRoom, err error){
//	result=new(ChatRoom)
//
//
//	var cName string
//	err=json2.Unmarshal(result.chatName,&cName);err==nil{
//		result.chatName=cName
//		return
//	}
//
//return
//}