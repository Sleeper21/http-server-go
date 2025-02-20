package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Sleeper21/http-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Create a struct that counts the number of requests received during a session
type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	secret         string
}

func main() {
	const port = "8080"

	godotenv.Load()              // Load environment variables
	dbURL := os.Getenv("DB_URL") // Assign them to a variable
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET must be set")
	}

	// Open connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	// Use your SQLC generated database package to create a new *database.Queries, and store it in your apiConfig struct so that handlers can access it:
	dbQueries := database.New(db)

	// Create an instance of apiConfig
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
		platform:       os.Getenv("PLATFORM"), // if platform != "dev" restrict dangerous endpoints
		secret:         os.Getenv("SECRET"),
	}

	// Create a new http.ServeMux
	mux := http.NewServeMux()

	// Use the Handle() method from the http.NewServeMux to add a handler for the "/" path
	// the "." in this case and we want to show the index.html file
	// By default by just using"." it will look for a file named index.html, otherwise it will show the list of all files in the directory"
	mux.Handle("/app/", apiCfg.addMetricsMiddleware(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	// Create handler functions to the endpoints
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerResetAccessToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToken)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgradeUserPlan)

	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	// Create a new http.Server struct
	// set Addr to ":8080"
	// Use the new ServeMux as the Handler
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	// Use the server's ListenAndServe method to start the server
	log.Printf("Listening on port %s", port)
	log.Fatal(server.ListenAndServe())

}
