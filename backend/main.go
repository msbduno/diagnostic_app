package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"diagnostic-backend/database"
	"diagnostic-backend/handlers"
)

const (
	defaultPort   = "8080"
	defaultDBPath = "./diagnostics.db"
)

func main() {
	log.Println("Démarrage du backend de diagnostic...")

	// Configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	//:= : Déclaration + assignation (type inféré)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// Initialiser la base de données
	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf(" Erreur d'initialisation de la base de données: %v", err)
	}
	defer database.CloseDB() //defer = exécute à la fin de la fonction main

	// Créer le routeur
	router := mux.NewRouter()

	// Routes API
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Diagnostics
	api.HandleFunc("/diagnostics", handlers.CreateDiagnostic).Methods("POST")
	api.HandleFunc("/diagnostics", handlers.GetDiagnostics).Methods("GET")
	api.HandleFunc("/diagnostics/{id:[0-9]+}", handlers.GetDiagnosticByID).Methods("GET")
	api.HandleFunc("/diagnostics/serial/{serial}", handlers.GetDiagnosticsBySerial).Methods("GET")

	// Statistiques
	api.HandleFunc("/statistics", handlers.GetStatistics).Methods("GET")

	// Middleware de logging
	//Un middleware est un intercepteur qui s'exécute avant chaque requête (comme un filtre en Java).
	router.Use(loggingMiddleware)

	// Configuration CORS
	//Sans CORS, le navigateur bloque les requêtes cross-origin qui permettent de communiquer entre le frontend et le backend.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // En production, spécifier les origines exactes
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	})

	handler := c.Handler(router)

	log.Printf("Serveur démarré sur le port %s", port)
	log.Printf("API disponible sur http://localhost:%s/api/v1", port)
	log.Printf(" Health check: http://localhost:%s/api/v1/health", port)
	log.Printf(" Base de données: %s", dbPath)
	log.Println(" Endpoints disponibles:")
	log.Println("   POST   /api/v1/diagnostics")
	log.Println("   GET    /api/v1/diagnostics")
	log.Println("   GET    /api/v1/diagnostics/{id}")
	log.Println("   GET    /api/v1/diagnostics/serial/{serial}")
	log.Println("   GET    /api/v1/statistics")
	log.Println("")

	// Démarrer le serveur
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf(" Erreur du serveur: %v", err)
	}
}

// loggingMiddleware enregistre toutes les requêtes
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log de la requête
		log.Printf("%s %s - %s", r.Method, r.RequestURI, r.RemoteAddr)

		// Passer à la prochaine étape
		next.ServeHTTP(w, r)

		// Log du temps de traitement
		duration := time.Since(start)
		log.Printf("Traité en %v", duration)
	})
}
