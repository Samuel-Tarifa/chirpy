package main

import "sync/atomic"

type apiConfig struct {
	fileserverHits atomic.Int32
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