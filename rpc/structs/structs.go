package structs

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
	Id      string
	Content string
}
type Messages struct {
	Messages []Message
}
type RequestCreateChatRoom struct {
	NameChatRoom string
}

type ResponseCreateChatRoom struct {
	Status string
}
type RequestListChatRoom struct {
	Username string
}
type ResponseListChatRoom struct {
	Status string
}
type RequestCreateClient struct {
	Username string
}
type ResponseCreateClient struct {
	Status string
	Client Client
}

type Client struct {
	Username string
	Id       string
}
type Clients struct {
	Clients []Client
}
