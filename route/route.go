package route

import (
	"encoding/json"
	"golang_socket/protocol"
	"log"
	"net"
	"os"
	"time"
)

type Msg struct {
	Conditions map[string]string `json:"meta"`
	Content    interface{}       `json:"content"`
}

type Response struct {
	StatusCode string      `json:"statusCode"`
	Result     interface{} `result`
}

func reader(conn net.Conn, readerChannel chan []byte, timeout int) {
	for {
		select {
		case data := <-readerChannel:
			conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			Business(conn, data)
			break
		case <-time.After(time.Second * time.Duration(timeout)):
			Log("It's really weird to get Nothing!!!")
			conn.Close()
			return
		}
	}
}
func Business(conn net.Conn, data []byte) {
	flag := false
	for _, v := range Routers {
		pred := v[0]
		act := v[1]
		var message Msg
		err := json.Unmarshal(data, &message)
		if err != nil {
			Log(err)
		}
		if pred.(func(entry Msg) bool)(message) {
			result := act.(Controller).Excute(message)
			_, err := WriteResult(conn, result)
			if err != nil {
				Log("conn.WriteResult()", err)
			}
			return
		}
	}
	if !flag {
		_, err := WriteError(conn, "1111", "不能处理此类型的业务")
		if err != nil {
			Log("conn.WriteError()", err)
		}
	}
}

func CheckError(err error) {
	if err != nil {
		log.Printf("Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

//长连接
func HandleConnection(conn net.Conn, timeout int) {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go reader(conn, readerChannel, timeout)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}

}

func WriteResult(conn net.Conn, result interface{}) (n int, err error) {
	data, err := json.Marshal(Response{StatusCode: "0000", Result: result})
	if err != nil {
		return 0, err
	}
	return conn.Write(data)
}

func WriteError(conn net.Conn, statusCode string, result interface{}) (n int, err error) {
	data, err := json.Marshal(Response{StatusCode: statusCode, Result: result})
	if err != nil {
		return 0, err
	}
	return conn.Write(data)
}

func Log(v ...interface{}) {
	log.Println(v...)
}

type Controller interface {
	Excute(message Msg) interface{}
}

var Routers [][2]interface{}

func Route(judge interface{}, controller Controller) {
	switch judge.(type) {
	case func(entry Msg) bool:
		{
			var arr [2]interface{}
			arr[0] = judge
			arr[1] = controller
			Routers = append(Routers, arr)
		}
	case map[string]string:
		{
			defaultJudge := func(entry Msg) bool {
				for keyjudge, valjudge := range judge.(map[string]string) {
					val, ok := entry.Conditions[keyjudge]
					if !ok {
						return false
					}
					if val != valjudge {
						return false
					}
				}
				return true
			}
			var arr [2]interface{}
			arr[0] = defaultJudge
			arr[1] = controller
			Routers = append(Routers, arr)
		}
	default:
		Log("Something is wrong in Router")
	}
}

type MirrorController struct {
}

func (this *MirrorController) Excute(message Msg) interface{} {
	_, err := json.Marshal(message)
	CheckError(err)
	if time.Now().Unix()%2 == 0 {
		return "失败"
	}
	return "消息推送成功"
}

// func mirrorHandle(entry Msg) bool {
// 	if entry.Conditions["msgtype"] == "binding" {
// 		return true
// 	}
// 	return false
// }

func init() {
	var mirror MirrorController
	Routers = make([][2]interface{}, 0, 10)
	kvs := make(map[string]string)
	kvs["msgtype"] = "binding"
	Route(kvs, &mirror)
}
