package blogServer

import(
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)


type Client struct{
	socket net.Conn
	data chan[]byte
}

type ClientManager struct{
	clients map[*Client] bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

func startServerMode(){
	fmt.Println("Starting server...")
	listener, error:=net.Listen("tcp",":12345")
	if error!=nil{
		fmt.Println(error)
	}
	manager:=ClientManager{
		clients: make(map[*Client] bool),
		broadcast:make(chan []byte),
		register:make(chan *Client),
		unregister:make(chan *Client),
	}
	go manager.start()
	for{
		connection,_:=listener.Accept()
		if error !=nil{
			fmt.Println(error)
		}
		//Here is adding the information of the connection to an instance of a client
		client:=&Client{socket: connection, data: make(chan []byte)}
		manager.register<-client
		go manager.receive(client)
		go manager.send(client)
	}

}
func startClientMode(){
	fmt.Println("Starting client...")
	connection, error:=net.Dial("tcp","localhost:12345")


	if error!=nil{
		fmt.Println(error)
	}
	client:=&Client{socket:connection}
	go client.receive()
	for{
		reader:=bufio.NewReader(os.Stdin)
		message,_:=reader.ReadString('\n')
		//TODO here i should create a username of the client fmt.Println(username);
		connection.Write([]byte(strings.TrimRight(message,"\n")))
	}

}


func (client *Client) print(){
	fmt.Println(client.data)
}

func(manager *ClientManager) start(){
	for{
		for k, v := range manager.clients {
			fmt.Println("k:", k, "v:", v)
			k.print()
		}
		select {
		case connection:= <-manager.register:
			manager.clients[connection]=true;
			fmt.Println("Added new connection")
		case connection:=<-manager.unregister:
			if _,ok:=manager.clients[connection];ok{
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("A connection has terminated")
			}
		case message:=<-manager.broadcast:
			for connection:=range manager.clients{
				select {
				case connection.data<-message:
				default:
					close(connection.data)
					delete(manager.clients,connection)
				}
			}
		}
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

func (manager *ClientManager) receive(client *Client){
	for{
		message:=make([]byte, 4096)
		length, err:=client.socket.Read(message)
		if err!=nil{
			manager.unregister<-client
			client.socket.Close()
			break
		}
		if length>0{
			fmt.Println("RECEIVED: "+string(message))
			manager.broadcast<-message
		}
	}
}

func (manager *ClientManager) send(client *Client){
	defer client.socket.Close()
	for{
		select{
		case message, ok:=<-client.data:
			if!ok{
				return
			}
			client.socket.Write(message)
		}
	}
}
func main(){
	//read the input, either client or server
	reader:=bufio.NewReader(os.Stdin);
	fmt.Print("Choose server('1') or client(2): ")
	//trusted command to read console
	input,_,err:=reader.ReadRune()
	if err!=nil{
		fmt.Println(err)
	}
	switch input {
	case '1':
		startServerMode()
	case '2':
		//options for client
		//fmt.Print()
		startClientMode()

	}
}