package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	log.Println("Listening at http://localhost:11112")
	http.ListenAndServe(":11111", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		b, err := json.Marshal(&answer{Immoral: 0.3, Person: 0.6, Obscene: 0.9})
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		log.Println(req.Model)
		w.WriteHeader(200)
		w.Write(b)
	}))
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
