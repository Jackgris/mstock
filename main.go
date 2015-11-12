package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/jackgris/mstock/controllers"
)

// This is gonna be our server
type Server struct {
	*http.ServeMux
}

// This server will only serve to respond to requests for static files,
// and we will add to our main server
var sm = http.NewServeMux()

// We will compile our templates one Time
var templates = template.Must(template.ParseFiles("templates/home.html"))

// We will verify the routes from which are going to be possible to
// access our server
var validPath = regexp.MustCompile("^/(auth)/([a-zA-Z0-9]+)$")

func main() {
	// only we create and launch the server
	server := NewServer()
	server.Run()
}

// NewServer it will create a server with all necessary
// settings to function properly
func NewServer() *Server {
	sm.Handle("/", http.FileServer(http.Dir("./public/")))
	s := new(Server)
	s.ServeMux = http.NewServeMux()
	s.HandleFunc("/auth/login", makeHandler(controllers.LoginHandler))
	s.HandleFunc("/auth/signup", makeHandler(controllers.RegisterHandler))
	s.HandleFunc("/", HomeHandler)
	return s
}

// Will run the server
func (s *Server) Run() {
	http.ListenAndServe(":8080", s)
}

// We use this function to verify that the route of the
// request is correct in each handler
func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

// It returns the main page.  Which will handle all the application
// view through angularjs
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// if the path corresponds to a static file simply delegate
	if strings.Contains(r.URL.Path, ".") {
		sm.ServeHTTP(w, r)
		return
	}
	// we show standard output data requests
	log.Println(r.Method, "path", r.URL.Path)
	// returns the compiled template to the user
	err := templates.ExecuteTemplate(w, "home.html", "")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
