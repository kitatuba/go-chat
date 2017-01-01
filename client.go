package main

import (
	"github.com/gorilla/websocket"
)

// clientはチャットを行っている1人のユーザを表現
type client struct {
	// socket 該当クライアントのためのwebSocket
	socket *websocket.Conn
	// send メッセージが送られるチャネル
	send chan []byte
	// 該当クライアントが参加しているチャットルーム
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err != nil {
			// 受け取ったメッセージは、forwardチャネルに送信される
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	// sendチャネルからメッセージを受け取り、書き出し
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
