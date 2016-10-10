# golang-socket
一个简单的golang socket服务框架，使用简单的通信协议解决粘包问题，通过心跳计时的方式能及时关闭长链接，自定义Route规则，调用Controller进行任务的分发处理
# 使用案例
## server
```
import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	"github.com/SmarkSeven/golang-socket/route"
)

type Controller struct {
}

func (this *Controller) Excute(message route.Message) interface{} {
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
	var controller Controller
	kvs := make(map[string]string)
 
	kvs["msgType"] = "send SMS"
  // 注册规则和处理器
	route.Route(kvs, &controller)
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
```

## client
```
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
   // 将数据打包后发送
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
```
