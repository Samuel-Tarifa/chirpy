package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Samuel-Tarifa/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	res := `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	w.Write([]byte(fmt.Sprintf(res, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "403 Forbidden")
		return
	}
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error deleting users:\n%v", err))
		return
	}
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Number of hits reseted."))
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	req := chirp{}
	err := decoder.Decode(&req)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("json invalid\n%v\n", err))
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, 400, "chirpy is too long")
		return
	}

	cleanBody := cleanString(req.Body)

	params := database.CreateChirpParams{
		Body:   cleanBody,
		UserID: req.User_id,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), params)

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error creating chirp:\n%v", err))
		return
	}

	resp := Chirp{
		ID:        chirp.ID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
		User_id:   chirp.UserID,
	}

	respondWithJSON(w, 201, resp)

}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Email string `json:"email"`
	}
	req := Req{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error decoding request:\n%v", err))
		return
	}

	u, err := cfg.db.CreateUser(r.Context(), req.Email)

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error creating user:\n%v", err))
		return
	}

	resp := User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
		Email:     u.Email,
	}

	respondWithJSON(w, 201, resp)

}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error getting chirps:\n%v", err))
		return
	}

	resp := []Chirp{}

	for _, chirp := range chirps {
		newChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt.Time,
			UpdatedAt: chirp.UpdatedAt.Time,
			Body:      chirp.Body,
			User_id:   chirp.UserID,
		}
		resp = append(resp, newChirp)
	}

	respondWithJSON(w, 200, resp)

}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("uuid invalid:\n%v", err))
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err == sql.ErrNoRows {
		respondWithError(w, 404, "no chirp found")
		return
	}
	if err != nil {
		respondWithError(w, 500, fmt.Sprint("error getting chirp:\n%v", err))
		return
	}
	resp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
		Body:      chirp.Body,
		User_id:   chirp.UserID,
	}
	respondWithJSON(w, 200, resp)
}
