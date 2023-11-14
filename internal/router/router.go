package router

import (
	"fmt"
	dbquery "hleb_flip/internal/db_query"
	"html/template"
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
	r.Router.HandleFunc("/player/{id:[0-9]+}", r.getPlayerHandler())
	r.Router.HandleFunc("/site", r.siteHandler())
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

		id := ro.Db.AddPlayer(b)
		log.Default().Printf("ok put request %s\n", string(b))

		io.WriteString(w, strconv.Itoa(id))
	}
}

func (ro *Router) changePlayerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Default().Println("handle change")
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

func (ro *Router) getPlayerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		log.Default().Printf("id requested %d\n", id)
		ans := ro.Db.GetPlayerRecord(id)
		io.WriteString(w, ans)
	}
}

func (ro *Router) siteHandler() http.HandlerFunc {
	tpl := template.Must(template.ParseFiles("index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		//io.WriteString(w, "<h1>Hello World!</h1>")
		tpl.Execute(w, nil)
	}
}
