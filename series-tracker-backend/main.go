// @title           Series Tracker API
// @version         1.0
// @description     API RESTful en Go para gestionar y rastrear el progreso de visualización de series de TV. Utiliza Chi para el enrutamiento y GORM para la base de datos.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api
// @schemes   http https

// Package main es el punto de entrada de la aplicación.
// Configura la base de datos, el router HTTP, el middleware, las rutas de la API,
// la ruta para Swagger UI y finalmente inicia el servidor web.
package main

import (
	"context" // Para el cierre grácil
	"log"
	"net/http"
	"os"
	"os/signal" // Para cierre grácil
	"syscall"   // Para cierre grácil
	"time"      // Para timeouts y cierre grácil

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"lab6/handlers"   // Asegúrate que la ruta de importación sea correcta
	"lab6/repository" // Asegúrate que la ruta de importación sea correcta

	// Importa http-swagger para servir la UI de Swagger
	httpSwagger "github.com/swaggo/http-swagger"

	// Importa los docs generados por swag init (IMPORTANTE el prefijo _ )
	_ "lab6/docs" // Asegúrate que la ruta de importación sea correcta (normalmente es "nombre_modulo/docs")
)

func main() {
	log.Println("Iniciando aplicación Series Tracker...")

	// Iniciar la conexión con la base de datos
	repository.InitDB()
	// Obtener la instancia de DB subyacente para poder cerrarla después
	sqlDB, err := repository.DB.DB()
	if err != nil {
		log.Fatalf("Error obteniendo la instancia DB subyacente: %v", err)
	}
	// Programar el cierre de la conexión DB al final de main
	defer func() {
		log.Println("Cerrando conexión a la base de datos...")
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error cerrando la conexión a la base de datos: %v", err)
		} else {
			log.Println("Conexión a la base de datos cerrada.")
		}
	}()
	// defer repository.CloseDB() // Alternativa si CloseDB maneja nil checks

	// Configurar router Chi
	r := chi.NewRouter()

	// --- Middleware ---
	// Middleware Logger: Loggea cada solicitud HTTP recibida (método, path, tiempo de respuesta)
	r.Use(middleware.Logger)
	// Middleware Recoverer: Recupera de panics, loggea el stack trace y devuelve un 500
	r.Use(middleware.Recoverer)
	// Middleware RequestID: Añade un ID único a cada solicitud para tracing
	r.Use(middleware.RequestID)
	// Middleware RealIP: Intenta obtener la IP real del cliente (detrás de proxies)
	r.Use(middleware.RealIP)

	// Middleware CORS: Configuración de Cross-Origin Resource Sharing
	// Esta configuración es permisiva, ¡restringir en producción!
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Permitir cualquier origen (inseguro para producción)
		// AllowedOrigins: []string{"http://localhost:3000", "https://mi-frontend.com"}, // Ejemplo más seguro
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Añadir headers necesarios
		ExposedHeaders:   []string{"Link"},                                                    // Headers expuestos al frontend
		AllowCredentials: true,                                                                // Permitir cookies/auth
		MaxAge:           300, // Tiempo máximo que el resultado de preflight puede ser cacheado (en segundos)
	}))

	// --- Rutas ---
	// Ruta para la documentación de Swagger UI
	// Servirá los archivos estáticos y el swagger.json generado
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	log.Println("Swagger UI disponible en /swagger/index.html")

	// Agrupar rutas de la API bajo el prefijo /api
	r.Route("/api", func(r chi.Router) {
		// Rutas para el recurso 'series'
		r.Get("/series", handlers.GetAllSeries)                // GET /api/series
		r.Post("/series", handlers.CreateSeries)               // POST /api/series
		r.Get("/series/{id}", handlers.GetSeriesByID)          // GET /api/series/123
		r.Put("/series/{id}", handlers.UpdateSeries)           // PUT /api/series/123
		r.Delete("/series/{id}", handlers.DeleteSeries)        // DELETE /api/series/123

		// Rutas de acciones específicas sobre 'series' (usando PATCH)
		r.Patch("/series/{id}/status", handlers.UpdateSeriesStatus)     // PATCH /api/series/123/status
		r.Patch("/series/{id}/episode", handlers.IncrementSeriesEpisode) // PATCH /api/series/123/episode
		r.Patch("/series/{id}/upvote", handlers.UpvoteSeries)           // PATCH /api/series/123/upvote
		r.Patch("/series/{id}/downvote", handlers.DownvoteSeries)       // PATCH /api/series/123/downvote
	})

	// Ruta de health check simple
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// --- Iniciar Servidor ---
	// Leer puerto desde variable de entorno o usar default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
	}
	serverAddr := ":" + port

	// Configurar servidor HTTP con timeouts y manejo de cierre grácil
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      r, // Usar el router Chi como handler
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Canal para escuchar errores del servidor en una goroutine separada
	serverErrors := make(chan error, 1)

	// Goroutine para iniciar el servidor
	go func() {
		log.Printf("Servidor escuchando en %s", serverAddr)
		serverErrors <- server.ListenAndServe()
	}()

	// --- Manejo de Cierre Grácil (Graceful Shutdown) ---
	// Canal para escuchar señales del sistema operativo (Interrupt, Terminate)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Bloquear hasta que se reciba un error del servidor o una señal de cierre
	select {
	case err := <-serverErrors:
		log.Fatalf("Error iniciando servidor: %v", err)

	case sig := <-shutdown:
		log.Printf("Señal de cierre (%v) recibida. Iniciando apagado grácil...", sig)

		// Crear un contexto con timeout para el apagado
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Intentar apagar el servidor grácilmente
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error durante el apagado grácil: %v", err)
			// Forzar cierre si Shutdown falla
			if err := server.Close(); err != nil {
				log.Printf("Error al forzar el cierre del servidor: %v", err)
			}
		} else {
			log.Println("Servidor apagado grácilmente.")
		}
	}

	log.Println("Aplicación terminada.")
}
