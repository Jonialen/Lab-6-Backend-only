package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"lab6/handlers"
	"lab6/repository"
)

func main() {
	// Iniciar la conexión con la base de datos
	repository.InitDB()

	// Configurar router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Considera restringir esto en producción
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Definir rutas para la API
	r.Get("/api/series", handlers.GetAllSeries)
	r.Get("/api/series/{id}", handlers.GetSeriesByID)
	r.Post("/api/series", handlers.CreateSeries)
	r.Put("/api/series/{id}", handlers.UpdateSeries)
	r.Delete("/api/series/{id}", handlers.DeleteSeries)

	r.Patch("/api/series/{id}/status", handlers.UpdateSeriesStatus)
	r.Patch("/api/series/{id}/episode", handlers.IncrementSeriesEpisode)
	r.Patch("/api/series/{id}/upvote", handlers.UpvoteSeries)
	r.Patch("/api/series/{id}/downvote", handlers.DownvoteSeries)

	// Iniciar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Servidor corriendo en http://localhost:" + port)
	http.ListenAndServe(":"+port, r)
}
