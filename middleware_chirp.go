package main

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/WhileCodingDoLearn/bootdev_server_tut/docs"
	"github.com/WhileCodingDoLearn/bootdev_server_tut/internal/database"
	"github.com/google/uuid"
)

type ChripMsg struct {
	Body string `json:"body"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

const UserKey = "user"

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:   "/api/chirps",
		Method: "POST",
		InputType: ChripMsg{
			Body: "Text",
		},
		OutputType:     Chirp{},
		Description:    "Endpoint for Sending Message to Chirpy",
		Authentication: true,
		ErrorCodes:     []int{http.StatusInternalServerError, http.StatusBadRequest, http.StatusCreated},
	})
}

func (apiCfg *apiConfig) ChirpHandler(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(UserKey).(User)

	data, errDecode := convertFromJSON(r.Body, ChripMsg{})
	if errDecode != nil {
		log.Printf("Error decoding parameter: %s", errDecode)
		respondWithError(w, http.StatusInternalServerError, "Error decoding request")
		return
	}

	if len(data.Body) > 140 || len(data.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	result := apiCfg.filter.FilterWord(data.Body)

	chirp, err := apiCfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   result,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, mapChirp(chirp))
}

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:           "/api/chirps/",
		Method:         "GET",
		QueryParams:    []string{"author_id=UUID", "sort=asc|desc"},
		InputType:      nil,
		OutputType:     []Chirp{Chirp{Body: "Text"}},
		Description:    "Endpoint for Getting Messages from User by Id",
		Authentication: false,
		ErrorCodes:     []int{http.StatusBadRequest, http.StatusOK, http.StatusInternalServerError},
	})
}

func (apiCfg *apiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")
	var data []database.Chirp
	if len(authorId) > 0 {

		idFromUser, errUuid := uuid.Parse(authorId)
		if errUuid != nil {
			responWithoutBody(w, http.StatusBadRequest)
			return
		}

		chirpsFromDB, errDB := apiCfg.dbQueries.GetChirpsFromUser(r.Context(), idFromUser)
		if errDB != nil {
			responWithoutBody(w, http.StatusBadRequest)
			return
		}
		data = chirpsFromDB
	} else {
		chirpsFromDB, err := apiCfg.dbQueries.GetAllShirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		data = chirpsFromDB
	}

	chirps := make([]Chirp, 0)
	for _, chirp := range data {
		chirps = append(chirps, mapChirp(chirp))
	}
	if strings.ToLower(sortOrder) == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:           "/api/chirps/{chirpID}",
		Method:         "GET",
		QueryParams:    []string{"chirpID=UUID"},
		InputType:      nil,
		OutputType:     Chirp{Body: "Text"},
		Description:    "Endpoint for Getting Messages by Chirp Id",
		Authentication: false,
		ErrorCodes:     []int{http.StatusBadRequest, http.StatusInternalServerError, http.StatusNotFound, http.StatusOK},
	})
}

func (apiCfg *apiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	errUuid := uuid.Validate(id)
	if errUuid != nil {
		respondWithError(w, http.StatusBadRequest, errUuid.Error())
		return
	}
	newUui, errParse := uuid.Parse(id)
	if errParse != nil {
		respondWithError(w, http.StatusInternalServerError, errParse.Error())
		return
	}
	data, errDB := apiCfg.dbQueries.GetChirpById(r.Context(), newUui)
	if errDB != nil {
		responWithoutBody(w, http.StatusNotFound)
		return
	}
	respondWithJSON(w, http.StatusOK, mapChirp(data))
}

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:           "/api/chirps/{chirpID}",
		Method:         "DELETE",
		QueryParams:    []string{"chirpID=UUID"},
		InputType:      nil,
		OutputType:     nil,
		Description:    "Endpoint for Deöeting Messages by Chirp Id",
		Authentication: true,
		ErrorCodes:     []int{http.StatusBadRequest, http.StatusInternalServerError, http.StatusNotFound, http.StatusForbidden, http.StatusNoContent},
	})
}

func (apiCfg *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")

	user := r.Context().Value(UserKey).(User)

	errUuid := uuid.Validate(id)
	if errUuid != nil {
		respondWithError(w, http.StatusBadRequest, errUuid.Error())
		return
	}
	chirpUui, errParse := uuid.Parse(id)
	if errParse != nil {
		respondWithError(w, http.StatusInternalServerError, errParse.Error())
		return
	}
	chirp, errDb := apiCfg.dbQueries.GetChirpById(r.Context(), chirpUui)
	if errDb != nil {
		respondWithError(w, http.StatusNotFound, errDb.Error())
		return
	}
	if chirp.UserID != user.ID {
		responWithoutBody(w, http.StatusForbidden)
		return
	}

	errDel := apiCfg.dbQueries.DeleteChirpById(r.Context(), chirpUui)
	if errDel != nil {
		respondWithError(w, http.StatusInternalServerError, errDel.Error())
		return
	}
	responWithoutBody(w, http.StatusNoContent)
}
