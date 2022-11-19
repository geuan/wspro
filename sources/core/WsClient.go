package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WsClient struct {
	conn *websocket.Conn
	readChan chan *WsMessage   // 读队列
	closeChan chan byte      // 失败队列
}

func NewWsClient(conn *websocket.Conn) *WsClient {
	return &WsClient{conn: conn,readChan:make(chan *WsMessage),closeChan:make(chan byte)}
}

func (this *WsClient) Ping(waitime time.Duration)  {
	for {
		time.Sleep(waitime)
		err := this.conn.WriteMessage(websocket.TextMessage,[]byte("ping"))
		if err != nil {
			ClientMap.Remove(this.conn)
			return
		}
	}
}

func (this *WsClient)  ReadLoop()  {
	for {
		t,data,err := this.conn.ReadMessage()
		if err!= nil {
			_ = this.conn.Close()
			ClientMap.Remove(this.conn)
			this.closeChan<-1
			break
		}
		this.readChan<-NewWsMessage(t,data)
	}

}

func (this *WsClient)  HandlerLoop() {
	loop:for {
		select {
		case msg:=<- this.readChan:
			fmt.Println(string(msg.MessageData))
		case <-this.closeChan:
			log.Println("已经关闭")
			break loop
		}
	}
}

















