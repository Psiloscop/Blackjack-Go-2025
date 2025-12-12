package contract

type Message struct {
	Text     string
	IsError  bool
	PlayerId PlayerId
	HandId   PlayerHandId
}

type MessageSender func(msg Message)
