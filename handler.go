package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/WhileCodingDoLearn/bootdev_server_tut/internal/database"
)

type Input struct {
	Body string `json:"body"`
}

type Output struct {
	Cleaned_Body string `json:"cleaned_body"`
}

type SuccessResponse struct {
	Valid bool `json:"valid"`
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte("{error:" + msg + "}"))
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	msg, errJSON := json.Marshal(payload)
	if errJSON != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error encoding response: %s", errJSON)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(msg)
}

func responWithoutBody(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func convertFromJSON[T interface{}](reader io.Reader, data T) (T, error) {
	var result T
	err := json.NewDecoder(reader).Decode(&data)
	if err != nil {
		return result, err
	}
	result = data
	return result, nil
}

func mapChirp(data database.Chirp) Chirp {
	return Chirp{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Body:      data.Body,
		UserID:    data.UserID,
	}
}
