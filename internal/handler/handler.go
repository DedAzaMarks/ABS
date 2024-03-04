package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	_ "github.com/valyala/fasthttp"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "postgres"
)

type Server struct {
	db *sql.DB
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	log.Printf("new request %s %s", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		switch string(r.URL.Path) {
		case "/ping":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("pong")); err != nil {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("can't respond on /ping request"))
			}
		default:
			log.Printf("%s not implemented", r.URL.Path)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("not implemeted"))
			return
		}
	case http.MethodPost:
	}
	log.Print("success")
}

func (s *Server) Run() error {
	return http.ListenAndServe(":42069", http.HandlerFunc(s.serve))
}

func NewServer() (*Server, error) {
	log.Print("connecting to db")
	dataSource := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", dataSource)
	defer func() { _ = db.Close() }()
	if err != nil {
		return nil, fmt.Errorf("can't open DB connection: %w", err)
	}
	return &Server{db: db}, nil
}
