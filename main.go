package main

import (
	"flag"
	"github.com/mnoster/chat/trace"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// templ represents a single template

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")

	flag.Parse()

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/room", r)
	http.HandleFunc("/auth/", loginHandler)

	// Startup room
	go r.run()

	log.Println("\n -- Web Server listening on port", *addr, "--")
	// Start Web Server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
