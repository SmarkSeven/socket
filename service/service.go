package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	"github.com/SmarkSeven/golang-socket/route"
)

type MirrorController struct {
}

func (this *MirrorController) Excute(message route.Message) interface{} {
	_, err := json.Marshal(message)
	CheckError(err)
	if time.Now().Unix()%2 == 0 {
		return "失败"
	}
	return "消息推送成功"
}

func CheckError(err error) {
	if err != nil {
		log.Printf("Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func Log(v ...interface{}) {
	log.Println(v...)
}

func init() {
	var mirror MirrorController
	kvs := make(map[string]string)
	kvs["msgType"] = "send SMS"
	route.Route(kvs, &mirror)
}

func main() {
	netListen, err := net.Listen("tcp", "localhost:6060")
	CheckError(err)
	defer netListen.Close()
	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		Log(conn.RemoteAddr().String(), " tcp connect success")
		// 如果此链接超过6秒没有发送新的数据，将被关闭
		go route.HandleConnection(conn, 6)
	}
}
