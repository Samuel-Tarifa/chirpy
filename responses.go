package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	dat := responseError{
		Error: msg,
	}
	res, err := json.Marshal(dat)
	if err != nil {
		log.Printf("error encoding json error response:\n%v", err)
		return
	}
	log.Print(msg)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	w.WriteHeader(code)
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshaling json response:\n%v", err)
		respondWithError(w, 500, "error")
	}
	w.WriteHeader(code)
	w.Write(res)
}
