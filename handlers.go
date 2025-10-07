package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Number of hits reseted."))
}

func validate_chirp(w http.ResponseWriter, r *http.Request) {

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

	dat := responseValid{
		Valid:        true,
		Cleaned_body: cleanBody,
	}
	respondWithJSON(w, 200, dat)

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
	}
	respondWithJSON(w, 200, u)

}
