package structs

import "time"

type OptionMessage struct {
	Option string
	Data   map[string]string
}

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

type Client struct {
	Username string
	Id       string
}
type Clients struct {
	Clients []Client
}
