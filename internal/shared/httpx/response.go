package httpx

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseHandler(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		log.Println("error encoding response:", err)
	}
}
