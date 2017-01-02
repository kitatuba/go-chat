package main

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
