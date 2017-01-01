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
		t.templ = template.Must(template.ParseFiles(filepath.Join("./../../templates", t.filename)))
	})

	t.templ.Execute(w, nil)
}

func main() {
	// ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})
	// WEBサーバの開始
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndserve:", err)
	}
}
