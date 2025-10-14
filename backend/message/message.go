package message

type MessageBody struct {
	Action string
	Arg    map[string]string
}
type Message struct {
	To  string
	Arg MessageBody
}
