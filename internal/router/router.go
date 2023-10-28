package router

import (
	"fmt"
	dbquery "hleb_flip/internal/db_query"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Router struct {
	Router *mux.Router
	Port   string
	Host   string
	Db     *dbquery.DB
}

func NewRouter(host, port string) *Router {
	return &Router{
		Router: mux.NewRouter(),
		Port:   port,
		Host:   host,
		Db:     dbquery.NewDB(),
	}
}

func (r *Router) StartRouter() {
	r.Router.HandleFunc("/hello", r.helloHandler())
	r.Router.HandleFunc("/db", r.dbHandler())
	r.Router.HandleFunc("/add", r.addPlayerHandler())
	r.Router.HandleFunc("/change", r.changePlayerHandler())
	http.ListenAndServe(fmt.Sprintf("%s:%s", r.Host, r.Port), r.Router)
	log.Default().Println("Server started")
}

func (r *Router) helloHandler() http.HandlerFunc {
	log.Default().Println("handling hello")
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
	}
}

func (ro *Router) dbHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := ro.Db.GetRecords()
		io.WriteString(w, a)
	}
}

func (ro *Router) addPlayerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := r.Header.Get("size")
		ken, err := strconv.Atoi(s)
		if err != nil {
			log.Default().Printf("f put request1 %+v\n", err)
		}
		b := make([]byte, ken)
		_, err = r.Body.Read(b)
		if err != nil && err != io.EOF {
			log.Default().Printf("f put request %+v\n", err)
		}

		ro.Db.AddPlayer(b)
		log.Default().Printf("ok put request %s\n", string(b))

		io.WriteString(w, string(b))
	}
}

func (ro *Router) changePlayerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := r.Header.Get("size")
		ken, err := strconv.Atoi(s)
		if err != nil {
			log.Default().Printf("f put request1 %+v\n", err)
		}
		b := make([]byte, ken)
		_, err = r.Body.Read(b)
		if err != nil && err != io.EOF {
			log.Default().Printf("f put request %+v\n", err)
		}

		ro.Db.ChangeRecordForPlayer(b)
		log.Default().Printf("ok put request %s\n", string(b))

		io.WriteString(w, string(b))
	}
}
