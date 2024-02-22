package main

import "net/http"

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	type Status struct {
		Ok bool
	}
	respondWithJSON(w, 200, Status{Ok: true})
}
