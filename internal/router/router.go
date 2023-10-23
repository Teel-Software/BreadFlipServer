package router

import (
	"fmt"
	dbquery "hleb_flip/internal/db_query"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
	port   string
	host   string
}

func NewRouter(host, port string) *Router {
	return &Router{
		router: mux.NewRouter(),
		port:   port,
		host:   host,
	}
}

func (r *Router) StartRouter() {
	r.router.HandleFunc("/hello", helloHandler())
	r.router.HandleFunc("/db", dbHandler())
	http.ListenAndServe(fmt.Sprintf("%s:%s", r.host, r.port), r.router)
	log.Default().Println("Server started")
}

func helloHandler() http.HandlerFunc {
	log.Default().Println("handling hello")
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
	}
}

func dbHandler() http.HandlerFunc {
	a := dbquery.Wtf()
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, a)
	}
}
