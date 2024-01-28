package router

import (
	"encoding/json"
	"fmt"
	dbquery "hleb_flip/internal/db_query"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Router struct {
	Router   *mux.Router
	Port     string
	Host     string
	Database *dbquery.DB
}

func NewRouter(host, port string) *Router {
	return &Router{
		Router:   mux.NewRouter(),
		Port:     port,
		Host:     host,
		Database: dbquery.NewDB(),
	}
}

func (ro *Router) StartRouter() {
	ro.Router.HandleFunc("/hello", ro.helloHandler()).Methods("GET", "OPTIONS", "POST", "PUT")
	ro.Router.HandleFunc("/db", ro.dbHandler()).Methods("GET", "OPTIONS", "POST", "PUT")
	ro.Router.HandleFunc("/getrecords", ro.getRecordsHandler()).Methods("GET", "OPTIONS", "POST", "PUT")
	ro.Router.HandleFunc("/add", ro.addPlayerHandler()).Methods("GET", "OPTIONS", "POST", "PUT")
	ro.Router.HandleFunc("/change", ro.changePlayerHandler()).Methods("GET", "OPTIONS", "POST", "PUT")
	ro.Router.HandleFunc("/player/{id:[0-9]+}", ro.getPlayerHandler()).Methods("GET", "OPTIONS", "POST", "PUT")
	ro.Router.HandleFunc("/site", ro.siteHandler()).Methods("GET", "OPTIONS", "POST", "PUT")
	http.ListenAndServe(fmt.Sprintf("%s:%s", ro.Host, ro.Port), ro.Router)
	log.Default().Println("Server started")
}

func (ro *Router) helloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		log.Default().Println("handling hello")
		io.WriteString(w, "hello")
	}
}

func (ro *Router) dbHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "OPTIONS":
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "*")
			w.Header().Add("Access-Control-Allow-Methods", "*")
		default:
			log.Default().Println("handling top ten")
			topRecords := ro.Database.GetTopTenRecords()
			io.WriteString(w, topRecords)
		}
	}
}

func (ro *Router) getRecordsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "OPTIONS":
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "*")
			w.Header().Add("Access-Control-Allow-Methods", "*")
		default:
			log.Default().Println("handling get page")
			offset, err := strconv.Atoi(r.Header.Get("offset"))
			if err != nil {
				log.Default().Printf("cant parse offset %+v\n", err)
			}
			count, err := strconv.Atoi(r.Header.Get("count"))
			if err != nil {
				log.Default().Printf("cant parse count %+v\n", err)
			}

			recordList := ro.Database.GetRecordsWithPaging(offset, count)

			marshalledRecords, err := json.Marshal(recordList)
			if err != nil {
				fmt.Fprintf(os.Stderr, "json failed: %v\n", err)
			}
			io.WriteString(w, string(marshalledRecords))
		}
	}
}

func (ro *Router) addPlayerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "OPTIONS":
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "*")
			w.Header().Add("Access-Control-Allow-Methods", "*")
		default:
			log.Default().Println("handling add player")
			newPlayerBytes := make([]byte, r.ContentLength)
			_, err := r.Body.Read(newPlayerBytes)
			if err != nil && err != io.EOF {
				log.Default().Printf("f put request %+v\n", err)
			}

			id := ro.Database.AddPlayer(newPlayerBytes)
			log.Default().Printf("ok put request %s\n", string(newPlayerBytes))

			io.WriteString(w, strconv.Itoa(id))
		}
	}
}

func (ro *Router) changePlayerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "OPTIONS":
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "*")
			w.Header().Add("Access-Control-Allow-Methods", "*")
		default:
			log.Default().Println("handling change player")
			reqBody := make([]byte, r.ContentLength)
			_, err := r.Body.Read(reqBody)
			if err != nil && err != io.EOF {
				log.Default().Printf("failed change player request %+v\n", err)
				return
			}

			ro.Database.ChangeRecordForPlayer(reqBody)
			log.Default().Printf("ok change player request %s\n", string(reqBody))
		}
	}
}

func (ro *Router) getPlayerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "OPTIONS":
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "*")
			w.Header().Add("Access-Control-Allow-Methods", "*")
		default:
			log.Default().Println("handling get player")
			id, _ := strconv.Atoi(mux.Vars(r)["id"])
			log.Default().Printf("id requested %d\n", id)
			ans := ro.Database.GetPlayerRecord(id)
			io.WriteString(w, ans)
		}
	}
}

func (ro *Router) siteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "OPTIONS":
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "*")
			w.Header().Add("Access-Control-Allow-Methods", "*")
		default:
			tpl := template.Must(template.ParseFiles("index.html"))
			log.Default().Println("handling site")
			tpl.Execute(w, ro.Database.GetRecordsWithPaging(0, 100).List)
		}
	}
}
