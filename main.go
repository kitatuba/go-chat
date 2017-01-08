package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// http.Handler型に適合させるため、ServeHTTPメソッドを実装
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("./templates", t.filename)))
	})

	t.templ.Execute(w, nil)
}

func main() {
	r := newRoom()
	// ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})
	// /room
	http.Handle("/room", r)
	// チャットルームを開始
	go r.run()
	// WEBサーバの開始
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndserve:", err)
	}
}
