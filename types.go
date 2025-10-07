package main

import (
	"sync/atomic"
	"time"

	"github.com/Samuel-Tarifa/chirpy/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}

type chirp struct {
	Body string `json:"body"`
}

type responseError struct {
	Error string `json:"error"`
}

type responseValid struct {
	Valid bool `json:"valid"`
	Cleaned_body string `json:"cleaned_body"`
}

type badWordsMap map[string]bool

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}