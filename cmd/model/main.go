package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	log.Println("Listening at http://localhost:11112")
	log.Println(http.ListenAndServe(":11112", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(&answer{Immoral: rand.Float64(), Person: rand.Float64(), Obscene: rand.Float64()})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println(req.Model)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(b); err != nil {
			log.Println(err)
		}
	})))
}

type request struct {
	Model string `json:"model"`
	Text  string `json:"text"`
}

type answer struct {
	Immoral float64 `json:"immoral"`
	Person  float64 `json:"person"`
	Obscene float64 `json:"obscene"`
}
