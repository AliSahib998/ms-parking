package util

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Encode(w http.ResponseWriter, v interface{}) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Error("error happened: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}
