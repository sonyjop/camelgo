package core

type Message struct {
	Body    interface{}
	Headers map[string]interface{}
}

type Exchange struct {
	ID         string
	InMessage  *Message
	OutMessage *Message
	Error      error
	Properties map[string]interface{}
}
