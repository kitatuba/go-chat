package main

type room struct {
	// forward 他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte
}
