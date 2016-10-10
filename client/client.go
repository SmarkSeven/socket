package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/SmarkSeven/golang-socket/protocol"

	"github.com/SmarkSeven/golang-socket/route"
)

type PushParam struct {
	CoachId     string `json:"coachId"`
	StudentName string `json:"studentName"`
	Phone       string `json:"phone"`
	// MsgType     string      `json:"msgType"`
	Datetime time.Time `json:"datetime"`
	// Extra    map[string]interface{} `json:"extras"`
}

// type Extras struct {
// 	StudentId int64
// }

type Info struct {
	Id       int64
	PushData string `json:"pushData"`
}
type Response struct {
	StatusCode string `json:"statusCode"`
	Result     string `json:"result"`
}

func senderMsg(conn net.Conn) {

	kvs := make(map[string]string)
	kvs["msgType"] = "send SMS"

	msg := route.Message{
		Conditions: kvs,
		Content: PushParam{
			CoachId:     "13",
			StudentName: "Sum",
			Phone:       "15108888888",
			Datetime:    time.Now(),
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Marchal err %#v", msg)
	}
	conn.Write(protocol.Packet(data))
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	var message Response
	err = json.Unmarshal(buffer[:n], &message)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s receive data string:%+v \n", conn.RemoteAddr().String(), message)

}

func main() {
	server := "localhost:6060"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	for i := 1; i < 100; i++ {
		senderMsg(conn)
		time.Sleep(time.Second)
	}

}
