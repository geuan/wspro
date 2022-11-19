package core

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
	"wspro/sources/models"
)

// 外部公共使用
var ClientMap *ClientMapStruct

func init() {
	ClientMap = &ClientMapStruct{}
}

type ClientMapStruct struct {
	data sync.Map // key是客户端IP value 就是 websocket连接对象
}

func (this *ClientMapStruct) Store(conn *websocket.Conn) {
	WsClient := NewWsClient(conn)
	this.data.Store(conn.RemoteAddr().String(), WsClient)
	go WsClient.Ping(time.Second * 1)  // 检查客户端是否正常
	go WsClient.ReadLoop()      // 处理读循环
	go WsClient.HandlerLoop()   // 处理总控制循环
}

// 向所有客户端发送消息
func (this *ClientMapStruct) SendAllPods() {
	this.data.Range(func(key, value interface{}) bool {
		c := value.(*WsClient).conn
		//err := c.WriteMessage(websocket.TextMessage, []byte(msg))
		err := c.WriteJSON(models.MockPodeList())
		if err != nil {
			this.Remove(c)
			log.Println(err)
		}
		return true
	})
}

func (this *ClientMapStruct) Remove(conn *websocket.Conn) {
	this.data.Delete(conn.RemoteAddr().String())
}
