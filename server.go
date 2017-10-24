package main

import (
	"fmt"
	"net"
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
		}
		messageResult:=make([]byte,4096)
		length,err:=conn.Read(messageResult)
		if length>0{
			//manage json
		}

	}
}