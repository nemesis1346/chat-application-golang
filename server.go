package main

import (
	"fmt"
	"net"
	json2 "encoding/json"
	"io"
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
		//messageResult:=make([]byte,4096)
		//length,err:=conn.Read(messageResult)
		//if length>0{
		//	//manage json
		//	fmt.Println(string(messageResult))
		//}

		//decode(conn)

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