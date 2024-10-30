package main

import(
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Client represents an individual client connection with a socket for communication
// and a data channel for sending messages to the client.
type Client struct{
	socket net.Conn
	data chan[]byte
}

// ClientManager manages the active clients, broadcasting, and client registrations.
type ClientManager struct{
	clients map[*Client] bool     // Active clients map with connected status.
	broadcast chan []byte         // Channel for broadcasting messages to clients.
	register chan *Client         // Channel for registering new clients.
	unregister chan *Client       // Channel for unregistering clients.
}

// startServerMode initializes and starts the server, setting up the listener
// and handling incoming connections by registering new clients.
func startServerMode(){
	fmt.Println("Starting server...")
	listener, error := net.Listen("tcp", ":12345")
	if error != nil {
		fmt.Println(error)
	}
	manager := ClientManager{
		clients:    make(map[*Client] bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()
	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		// Add connection info to a new client instance
		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}

// startClientMode initializes and starts the client, connecting to the server,
// and enabling message input and sending.
func startClientMode(){
	fmt.Println("Starting client...")
	connection, error := net.Dial("tcp", "localhost:12345")
	if error != nil {
		fmt.Println(error)
	}
	client := &Client{socket: connection}
	go client.receive()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		// Future implementation: add client username
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
}

// print outputs the data channel information of the client.
func (client *Client) print(){
	fmt.Println(client.data)
}

// start launches the client manager, managing client connections and handling
// broadcast messaging.
func (manager *ClientManager) start(){
	for {
		for k, v := range manager.clients {
			fmt.Println("k:", k, "v:", v)
			k.print()
		}
		select {
		case connection := <- manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection")
		case connection := <- manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("A connection has terminated")
			}
		case message := <- manager.broadcast:
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

// receive listens for incoming messages on the client's socket,
// and handles message display or client disconnection.
func (client *Client) receive(){
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message))
		}
	}
}

// receive listens for incoming messages on a specific client connection,
// broadcasting received messages to all other connected clients.
func (manager *ClientManager) receive(client *Client){
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(message))
			manager.broadcast <- message
		}
	}
}

// send continuously sends messages from the data channel to the client socket,
// ensuring the connection is properly closed when no more messages are available.
func (manager *ClientManager) send(client *Client){
	defer client.socket.Close()
	for {
		select {
		case message, ok := <- client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}

// main function to start the program in either server or client mode,
// based on user input.
func main(){
	// Read input to choose server or client mode
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose server('1') or client(2): ")
	input, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	switch input {
	case '1':
		startServerMode()
	case '2':
		startClientMode()
	}
}
