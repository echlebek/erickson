package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func jsonError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{err.Error()})
	if err != nil {
		log.Println(err)
	}
}
