package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	c "github.com/jcardenasc93/work-at-olist/app/controllers"
	m "github.com/jcardenasc93/work-at-olist/app/middlewares"
	"github.com/jcardenasc93/work-at-olist/app/models"
)

type APIServer struct {
	port       string
	production bool
}

func NewAPIServer(port string, production bool) *APIServer {
	return &APIServer{
		port:       port,
		production: production,
	}
}

func (s *APIServer) Run() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/authors", func(r chi.Router) {
		r.With(m.Pagination).Get("/", c.HTTPHandleFunc(c.GetAuthors))
	})

	log.Printf("Server active on port: %s", s.port)
	log.Printf("Production: %v", s.production)
	http.ListenAndServe(s.port, r)

}

func main() {
	models.InitDB()
	server := NewAPIServer(":8080", false)
	server.Run()

}
