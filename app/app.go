package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jcardenasc93/work-at-olist/app/controllers"
	"github.com/jcardenasc93/work-at-olist/app/models"
)

func main() {
	models.InitDB()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/authors", func(r chi.Router) {
		r.Get("/", controllers.GetAuthors)
	})

	log.Println("Server active on port:", 8080)
	http.ListenAndServe(":8080", r)
}
