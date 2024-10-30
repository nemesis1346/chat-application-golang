# Blog Server Project

## Project Overview

This project implements a basic TCP-based chat server and client application in Go. It enables multiple clients to connect to a server, send messages, and receive broadcasts from other clients in real time. This code demonstrates the use of sockets for network communication and basic concurrency patterns in Go, with channels managing message passing and client interactions.

## Features

- **Server Mode**
  - Listens for incoming client connections on a specific port.
  - Manages clients using a `ClientManager`, which tracks active clients.
  - Broadcasts messages from any client to all other connected clients.
  - Handles client disconnections and automatically unregisters them.
  
- **Client Mode**
  - Connects to the server and continuously reads user input from the console.
  - Sends messages to the server for broadcast to other clients.
  - Displays broadcast messages from other clients in real time.
  
- **Concurrency Management**
  - Server handles each client in a separate goroutine for concurrent message processing.
  - Uses channels for broadcasting messages and managing client registration/unregistration.

## Code Structure

### 1. `Client` Struct

Represents an individual client connection, with:
- A `socket` for sending and receiving messages.
- A `data` channel to queue messages for outgoing data.

### 2. `ClientManager` Struct

The `ClientManager` struct manages all active client connections:
- `clients` map for tracking active clients.
- `broadcast` channel for distributing messages to clients.
- `register` and `unregister` channels for adding or removing clients.

### 3. `startServerMode` Function

Initializes the server, setting up a listener on port 12345:
- Accepts new client connections, creating a `Client` instance for each one.
- Registers each client with `ClientManager`.
- Starts two goroutines per client for concurrent message receiving and sending.

### 4. `startClientMode` Function

Initiates a client connection to the server:
- Prompts the user for console messages.
- Sends messages to the server, which broadcasts them to all clients.
- Displays incoming messages from other clients.

### 5. `ClientManager` Methods

The `ClientManager` methods include:
- `start`: Manages broadcasts and client registration/unregistration.
- `receive`: Listens for messages from a specific client and broadcasts them.
- `send`: Sends messages to a client from the `data` channel.

### 6. `main` Function

The `main` function prompts the user to select server or client mode:
- Server mode starts listening for client connections.
- Client mode connects to the server and allows message input.

## Getting Started

### Prerequisites
- [Go](https://golang.org/doc/install) installed.

### Running the Server

1. Run the following command:
   ```bash
   go run main.go
