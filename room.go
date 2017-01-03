package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// プロジェクト内で何度もハードコートされる可能性がある値については、定数として宣言
const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// upgrader HTTP接続をアップグレードする
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

// ServeHTTP http.Handler型に適合 HTTPハンドラとして扱えるようになる
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// WebSocketコネクションを取得
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	// チャットクライアント(ユーザ)を生成
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	// 該当のチャットルームのjoinチャネルにクライアントを送信
	r.join <- client

	// クライアントの終了時に退室処理を必ず行う
	defer func() { r.leave <- client }()

	// 別goroutineに担当させる
	go client.write()

	// sendチャネルのクローズまで接続を保持
	client.read()
}

type room struct {
	// forward 他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte
	// join チャットルームに参加しようとしてるクライアントのためのチャネル
	join chan *client
	// leave チャットルームから退室しようとしているクライアントのためのチャネル
	leave chan *client
	// clients 在室しているすべてのクライアントを保持するマップ join及びleaveプロパティを通して操作される
	clients map[*client]bool
}

func (r *room) run() {

	// 無限ループ 強制終了まで実行を継続
	for {
		select {
		case client := <-r.join:
			// 参加 clientsは、スライスでも代用出来るが入退室の繰り返しで無駄な要素が増えるためメモリ容量対策にマップを適用
			r.clients[client] = true
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// 全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
				// メッセージを送信
				default:
					// 送信に失敗
					delete(r.clients, client)
					// client.writeメソッドのforループを終了させる
					close(client.send)
				}
			}
		}
	}
}
