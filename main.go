package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/mnoster/chat/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/subosito/gotenv"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, r)
}

func main() {
	gotenv.Load()

	var addr = flag.String("addr", ":8080", "The addr of the application.")

	flag.Parse()

	gomniauth.SetSecurityKey("x1x2x3x4x5x6x7x8x9x0")

	gomniauth.WithProviders(
		facebook.New("key", "secret", "http://localhost:8080/auth/callback/facebook"),

		google.New(os.Getenv("GOOGLE_URL"), os.Getenv("GOOGLE_SECRET"), "http://localhost:8080/auth/callback/google"),

		github.New("key", "secret", "http://localhost:8080/auth/callback/github"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	// MustAuth referenced from auth.go file
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
