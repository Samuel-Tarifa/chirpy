package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Samuel-Tarifa/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	tokenSecret := os.Getenv("JWT_SECRET")

	apiCfg := apiConfig{}
	apiCfg.platform = os.Getenv("PLATFORM")
	apiCfg.polkaKey = os.Getenv("POLKA_KEY")
	apiCfg.tokenSecret = tokenSecret

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	apiCfg.db = dbQueries

	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(makeFileServer(".")))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK\n"))
	})

	mux.HandleFunc("GET /admin/metrics", apiCfg.getHits)

	mux.HandleFunc("POST /admin/reset", apiCfg.resetHits)

	mux.HandleFunc("POST /api/users", apiCfg.createUser)

	mux.HandleFunc("PUT /api/users", apiCfg.updateUser)

	mux.HandleFunc("POST /api/login", apiCfg.loginUser)

	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)

	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)

	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)

	mux.HandleFunc("POST /api/refresh", apiCfg.refreshAccessToken)

	mux.HandleFunc("POST /api/revoke", apiCfg.revokeRefreshToken)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.upgradeUser)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("error starting the server:\n%v", err)
	}
}
