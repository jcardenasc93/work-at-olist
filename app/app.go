package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	c "github.com/jcardenasc93/work-at-olist/app/controllers"
	"github.com/jcardenasc93/work-at-olist/app/db"
	m "github.com/jcardenasc93/work-at-olist/app/middlewares"
)

type APIServer struct {
	port       string
	production bool
	db         db.ApiDB
}

func NewAPIServer(port string, production bool, db db.ApiDB) *APIServer {
	return &APIServer{
		port:       port,
		production: production,
		db:         db,
	}
}

func (s *APIServer) Run() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/authors", func(r chi.Router) {
		r.With(m.Pagination).Get("/", c.HTTPHandleFunc(c.GetAuthors, s.db))
		r.Post("/", c.HTTPHandleFunc(c.CreateBook, s.db))
	})

	log.Printf("Server active on port: %s", s.port)
	log.Printf("Production: %v", s.production)
	http.ListenAndServe(s.port, r)

}

func main() {
	db, err := db.NewSQLiteDB()
	if err != nil {
		log.Fatal("Couldn't initialize DB")
	}
	err = db.Setup()
	if err != nil {
		log.Fatal("Couldn't initialize DB")
	}
	server := NewAPIServer(":8080", false, db)
	server.Run()
}
