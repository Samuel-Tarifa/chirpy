package main

import (
	"sync/atomic"
	"time"

	"github.com/Samuel-Tarifa/chirpy/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

type chirp struct {
	Body    string    `json:"body"`
	User_id uuid.UUID `json:"user_id"`
}

type responseError struct {
	Error string `json:"error"`
}

type badWordsMap map[string]bool

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	User_id   uuid.UUID `json:"user_id"`
}
