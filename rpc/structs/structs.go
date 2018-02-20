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
	Status   string
	ChatRoom ChatRoom
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
type ResquestGetClient struct {
	Username string
}
type ResponseGetClient struct {
	Client Client
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

type Client struct {
	Username string
	Id       string
}
type Clients struct {
	Clients []Client
}
