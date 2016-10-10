package socket

type Controller interface {
	Excute(message Message) interface{}
}
