package httpresponse

import (
	"encoding/json"
	"log"
	"net/http"
)

// Send func
func Send(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	resp, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(resp)
}
