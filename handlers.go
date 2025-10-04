package main

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Printf("error decoding json:\n%v", err)
		dat := responseError{
			Error: "Json invalid",
		}
		res, err := json.Marshal(dat)
		if err != nil {
			log.Printf("error encoding json error response:\n%v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
		w.WriteHeader(400)
		return
	}

	if len(req.Body)>140{
		dat := responseError{
			Error: "Chirp is too long",
		}
		res, err := json.Marshal(dat)
		if err != nil {
			log.Printf("error encoding json error response:\n%v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
		w.WriteHeader(400)
		return
	}

	cleanBody:=cleanString(req.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	dat := responseValid{
		Valid: true,
		Cleaned_body:cleanBody,
	}
	res,err:=json.Marshal(dat)
	if err!=nil{
			log.Printf("error encoding json response:\n%v", err)
			return
	}
	w.Write(res)

}
