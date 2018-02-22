package structs

import "time"

type ChatRooms struct {
	Chats []ChatRoom
}

type ChatRoom struct {
	NameChatRoom string
	Clients      Clients
	Id           string
	Messages     Messages
}
type Message struct {
	Id           string
	Content      string
	Username     string
	NameChatRoom string
	Time         time.Time
}
type Messages struct {
	Messages []Message
}
type RequestCreateChatRoom struct {
	NameChatRoom string
}

type ResponseCreateChatRoom struct {
	Status   string
	ChatRoom ChatRoom
}
type RequestListChatRoom struct {
	Username string
}
type ResponseListChatRoom struct {
	Status    string
	ChatRooms ChatRooms
}
type RequestCreateClient struct {
	Username string
}
type ResponseCreateClient struct {
	Status string
	Client Client
}
type RequestGetClient struct {
	Username string
}
type ResponseGetClient struct {
	Client Client
	Status string
}
type RequestGetChatRoom struct {
	ChatRoomName string
}
type ResponseGetChatRoom struct {
	ChatRoom ChatRoom
	Status   string
}
type RequestJoinChatRoom struct {
	Client   Client
	ChatRoom ChatRoom
}
type ResponseJoinChatRoom struct {
	Client Client
	Status string
}
type RequestLeaveChatRoom struct {
	Client   Client
	ChatRoom ChatRoom
}
type ResponseLeaveChatRoom struct {
	Client   Client
	ChatRoom ChatRoom
	Status   string
}
type RequestSaveMessage struct {
	Client   Client
	Content  string
	ChatRoom ChatRoom
	Time     time.Time
}
type ResponseSaveMessage struct {
	Status   string
	Messages Messages
}
type RequestGetPreviousMessages struct {
	ChatRoom ChatRoom
	Client   Client
}
type ResponseGetPreviousMessages struct {
	Messages Messages
	Status   string
}

type Client struct {
	Username string
	Id       string
}
type Clients struct {
	Clients []Client
}
