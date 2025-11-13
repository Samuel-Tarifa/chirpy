package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Samuel-Tarifa/chirpy/internal/auth"
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

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	id, err := auth.ValidateJWT(bearer, cfg.tokenSecret)

	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	req := chirp{}
	err = decoder.Decode(&req)

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
		UserID: id,
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := Req{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error decoding request:\n%v", err))
		return
	}

	if req.Email == "" || req.Password == "" {
		respondWithError(w, 400, "need to add email and password to request")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, 500, "error hashing password")
	}
	params := database.CreateUserParams{
		Email: req.Email,
		HashedPassword: sql.NullString{
			String: hash,
			Valid:  true,
		},
	}

	u, err := cfg.db.CreateUser(r.Context(), params)

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error creating user:\n%v", err))
		return
	}

	resp := UserResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
		Email:     u.Email,
		IsChirpyRed: u.IsChirpyRed.Bool,
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
		respondWithError(w, 500, fmt.Sprintf("error getting chirp:\n%v", err))
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

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := Req{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("error decoding request:\n%v", err))
		return
	}

	if req.Email == "" || req.Password == "" {
		respondWithError(w, 400, "need to add email and password to request")
		return
	}

	u, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, 404, "invalid credentials")
		return
	}
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error getting user %v", err))
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, u.HashedPassword.String)

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	if !match {
		respondWithError(w, 401, "401 Unauthorized")
		return
	}

	token, err := auth.MakeJWT(u.ID, cfg.tokenSecret, time.Hour)

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	params := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}

	createdRefreshToken, err := cfg.db.CreateRefreshToken(r.Context(), params)

	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	resp := UserResponseWithToken{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
		Email:        u.Email,
		Token:        token,
		RefreshToken: createdRefreshToken.Token,
		IsChirpyRed: u.IsChirpyRed.Bool,
	}

	respondWithJSON(w, 200, resp)

}

func (cfg *apiConfig) refreshAccessToken(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
	}
	token, err := cfg.db.GetRefreshToken(r.Context(), bearer)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 401, "Invalid token")
			return
		}
		respondWithError(w, 500, "error getting token")
		return
	}

	if token.ExpiresAt.Before(time.Now()) || token.RevokedAt.Valid {
		respondWithError(w, 401, "token expired")
		return
	}

	accessToken, err := auth.MakeJWT(token.UserID, cfg.tokenSecret, time.Hour)

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error creating access token %v", err))
	}

	res := TokenResp{Token: accessToken}

	respondWithJSON(w, 200, res)
}

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
	}
	params := database.UpdateRefreshTokenParams{
		Token:     bearer,
		ExpiresAt: time.Now(),
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	_, err = cfg.db.UpdateRefreshToken(r.Context(), params)

	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	respondWithJSON(w, 204, nil)
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(bearer, cfg.tokenSecret)

	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	v := req{}

	if err = decoder.Decode(&v); err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	if v.Email == "" {
		respondWithError(w, 400, "invalid email")
		return
	}
	if v.Password == "" {
		respondWithError(w, 400, "invalid password")
		return
	}

	hashedPassword, err := auth.HashPassword(v.Password)

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	params := database.UpdateUserParams{
		ID:    userID,
		Email: v.Email,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	u, err := cfg.db.UpdateUser(r.Context(), params)

	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	resp := UserResponse{
		ID:        u.ID,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
		Email:     u.Email,
		IsChirpyRed: u.IsChirpyRed.Bool,
	}

	respondWithJSON(w, 200, resp)
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {

	bearer, err := auth.GetBearerToken(r.Header)

	userIDFromToken, err := auth.ValidateJWT(bearer, cfg.tokenSecret)

	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	chirpIDUrl := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDUrl)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, err.Error())
			return
		}
		respondWithError(w, 500, err.Error())
		return
	}

	if chirp.UserID != userIDFromToken {
		respondWithError(w, 403, "not authorized")
		return
	}

	_, err = cfg.db.DeleteChirp(r.Context(), chirpID)

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 204, nil)

}

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	type user struct {
		UserID string `json:"user_id"`
	}
	type req struct {
		Event string `json:"event"`
		Data  user   `json:"data"`
	}

	data := req{}

	if err:=decoder.Decode(&data);err!=nil{
		respondWithError(w,500,err.Error())
	}

	if data.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}

	uuid,err:=uuid.Parse(data.Data.UserID)

	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("invalid uuid: %s",err))
		return
	}

	params:=database.SetChirpyRedUserParams{
		ID: uuid,
		IsChirpyRed: sql.NullBool{
			Bool: true,
			Valid: true,
		},
	}

	_,err=cfg.db.SetChirpyRedUser(r.Context(),params)

	if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			respondWithError(w,404,"user not found")
			return
		}
		respondWithError(w,500,err.Error())
	}

	respondWithJSON(w,204,nil)
}
