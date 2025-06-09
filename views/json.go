package views

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseId struct {
	ID int `json:"id"`
}

type ResponseIdStr struct {
	ID string `json:"id"`
}

func NewResponseId(id int) ResponseId {
	return ResponseId{
		ID: id,
	}
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	RespondWithJSON(w, code, ErrorResponse{
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Message is set: %s", err)
	}
}
