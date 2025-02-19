package main

import (
	"context"
	"net/http"
	"time"

	"github.com/WhileCodingDoLearn/bootdev_server_tut/docs"
	"github.com/WhileCodingDoLearn/bootdev_server_tut/internal/auth"
	"github.com/google/uuid"
)

type ResponseToken struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (apiCfg *apiConfig) isAutneticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt, errJwt := auth.GetToken(r.Header, "Bearer")
		if errJwt != nil {
			responWithoutBody(w, http.StatusUnauthorized)
			return
		}

		userId, errValidate := auth.ValidateJWT(jwt, apiCfg.secret)
		if errValidate != nil {
			responWithoutBody(w, http.StatusUnauthorized)
			return
		}
		user, errUserNotFound := apiCfg.dbQueries.GetUserByID(r.Context(), userId)
		if errUserNotFound != nil {
			responWithoutBody(w, http.StatusUnauthorized)
			return
		}
		ctxWithUser := context.WithValue(r.Context(), "user",
			User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, IsChirpyRed: user.IsChirpyRed})
		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}

type TokenResponse struct {
	Token string `json:"token"`
}

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:           "/api/refresh",
		Method:         "POST",
		InputType:      nil,
		OutputType:     TokenResponse{},
		Description:    "Enpoint to refresh Access Token for a User",
		Authentication: true,
		ErrorCodes:     []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusOK},
	})
}

func (apiCfg *apiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {

	token, errToken := auth.GetToken(r.Header, "Bearer")

	if errToken != nil {
		responWithoutBody(w, http.StatusUnauthorized)
		return
	}

	refreshTokenFromDB, errInDB := apiCfg.dbQueries.GetRefreshToken(r.Context(), token)
	if errInDB != nil {
		responWithoutBody(w, http.StatusUnauthorized)
		return
	}

	if refreshTokenFromDB.RevokedAt.Valid {
		responWithoutBody(w, http.StatusUnauthorized)
		return
	}

	if refreshTokenFromDB.ExpiresAt.Before(time.Now()) {
		responWithoutBody(w, http.StatusUnauthorized)
		return
	}

	jwt, errJWT := auth.MakeJWT(refreshTokenFromDB.UserID, apiCfg.secret, apiCfg.validTokenTime)
	if errJWT != nil {
		respondWithError(w, http.StatusInternalServerError, errJWT.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, TokenResponse{
		Token: jwt,
	})
}

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:           "/api/revoke",
		Method:         "POST",
		InputType:      nil,
		OutputType:     nil,
		Description:    "Enpoint to revoke provided refresh Token for a User",
		Authentication: true,
		ErrorCodes:     []int{http.StatusUnauthorized, http.StatusInternalServerError, http.StatusNoContent},
	})
}

func (apiCfg *apiConfig) RevokeToken(w http.ResponseWriter, r *http.Request) {
	token, errToken := auth.GetToken(r.Header, "Bearer")

	if errToken != nil {
		responWithoutBody(w, http.StatusUnauthorized)
		return
	}
	errRevoke := apiCfg.dbQueries.RevokeRefreshToken(r.Context(), token)
	if errRevoke != nil {
		responWithoutBody(w, http.StatusInternalServerError)
		return
	}
	responWithoutBody(w, http.StatusNoContent)
}
