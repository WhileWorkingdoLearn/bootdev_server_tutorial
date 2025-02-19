package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"time"

	"github.com/WhileCodingDoLearn/bootdev_server_tut/docs"
	"github.com/WhileCodingDoLearn/bootdev_server_tut/internal/auth"
	"github.com/WhileCodingDoLearn/bootdev_server_tut/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}
type UserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func init() {
	var Handler = docs.GetDocsGenerator()
	Handler.AddDoc(docs.EndpointDescription{
		Name:           "/api/users",
		Method:         "POST",
		InputType:      UserRequest{},
		OutputType:     User{},
		Description:    "Enpoint to register a User",
		Authentication: false,
		ErrorCodes:     []int{http.StatusBadRequest, http.StatusInternalServerError},
	})
}

func (apiCfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	errDecode := json.NewDecoder(r.Body).Decode(&req)
	if errDecode != nil {
		respondWithError(w, http.StatusBadRequest, errDecode.Error())
		return
	}

	_, errEmail := mail.ParseAddress(req.Email)
	if errEmail != nil {
		respondWithError(w, http.StatusBadRequest, errEmail.Error())
		return
	}

	_, errMailNotInDB := apiCfg.dbQueries.GetUserByEmail(r.Context(), req.Email)
	if errMailNotInDB == nil {
		respondWithError(w, http.StatusBadRequest, "Email already taken")
		return
	}

	pwhash, errPw := auth.HashPassword(req.Password)
	if errPw != nil {
		respondWithError(w, http.StatusBadRequest, errPw.Error())
		return
	}
	user, errDB := apiCfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: pwhash,
	})
	if errDB != nil {
		respondWithError(w, http.StatusInternalServerError, errDB.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})

}

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:           "/api/login",
		Method:         "POST",
		InputType:      UserRequest{},
		OutputType:     ResponseToken{},
		Description:    "Entpoint for User to log in",
		Authentication: false,
		ErrorCodes:     []int{400, 401, 500, 200},
	})
}

func (apiCfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	var requestFormat UserRequest
	errDecode := json.NewDecoder(r.Body).Decode(&requestFormat)
	if errDecode != nil {
		respondWithError(w, http.StatusBadRequest, errDecode.Error())
		return
	}

	_, errEmail := mail.ParseAddress(requestFormat.Email)
	if errEmail != nil {
		respondWithError(w, http.StatusBadRequest, errEmail.Error())
		return
	}

	user, errMailNotInDB := apiCfg.dbQueries.GetUserByEmail(r.Context(), requestFormat.Email)
	if errMailNotInDB != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	errPw := auth.CheckPasswordHash(requestFormat.Password, user.HashedPassword)
	if errPw != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password or password")
		return
	}

	accesstoken, errToken := auth.MakeJWT(user.ID, apiCfg.secret, apiCfg.validTokenTime)
	if errToken != nil {
		respondWithError(w, http.StatusInternalServerError, errToken.Error())
		return
	}

	refreshToken, errRefToken := auth.MakeRefreshToken()
	if errRefToken != nil {
		respondWithError(w, http.StatusInternalServerError, errRefToken.Error())
		return
	}

	_, errRefTokenDB := apiCfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if errRefToken != nil {
		respondWithError(w, http.StatusInternalServerError, errRefTokenDB.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ResponseToken{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        accesstoken,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	})
}

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:           "/api/users",
		Method:         "PUT",
		InputType:      UserRequest{},
		OutputType:     User{},
		Description:    "Entpoint for updating Users email and password",
		Authentication: true,
		ErrorCodes:     []int{http.StatusBadRequest, http.StatusInternalServerError, http.StatusOK},
	})
}

func (apiCfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(User)
	var req UserRequest
	errDecode := json.NewDecoder(r.Body).Decode(&req)
	if errDecode != nil {
		respondWithError(w, http.StatusBadRequest, errDecode.Error())
		return
	}

	_, errEmail := mail.ParseAddress(req.Email)
	if errEmail != nil {
		respondWithError(w, http.StatusBadRequest, errEmail.Error())
		return
	}

	pwhash, errPw := auth.HashPassword(req.Password)
	if errPw != nil {
		respondWithError(w, http.StatusBadRequest, errPw.Error())
		return
	}

	errUpdate := apiCfg.dbQueries.UpdateEmailAndPassword(r.Context(), database.UpdateEmailAndPasswordParams{
		ID:             user.ID,
		Email:          req.Email,
		HashedPassword: pwhash,
	})
	if errUpdate != nil {
		respondWithError(w, http.StatusInternalServerError, errUpdate.Error())
		return
	}

	user.Email = req.Email

	respondWithJSON(w, http.StatusOK, user)
}

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:           "/api/users",
		Method:         "DELETE",
		InputType:      nil,
		OutputType:     nil,
		Description:    "Entpoint for Deleting a User by its id",
		Authentication: true,
		ErrorCodes:     []int{http.StatusNotFound, http.StatusNoContent},
	})
}

func (apiCfg *apiConfig) deleteUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)
	err := apiCfg.dbQueries.DeleteUserById(r.Context(), user.ID)
	if err != nil {
		responWithoutBody(w, http.StatusNotFound)
		return
	}
	responWithoutBody(w, http.StatusNoContent)
}

type Event struct {
	Event string `json:"event"`
	Data  struct {
		User_id uuid.UUID `json:"user_Id"`
	} `json:"data"`
}

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:           "/api/polka/webhooks",
		Method:         "POST",
		InputType:      Event{},
		OutputType:     nil,
		Description:    "Webhook for changing Users subscription status",
		Authentication: true,
		ErrorCodes:     []int{http.StatusBadRequest, http.StatusInternalServerError, http.StatusOK},
	})
}

func (apiConfig *apiConfig) handleEvent(w http.ResponseWriter, r *http.Request) {
	reqApiKey, errKey := auth.GetToken(r.Header, "ApiKey")

	if errKey != nil || len(apiConfig.webhookKey) == 0 || reqApiKey != apiConfig.webhookKey {
		responWithoutBody(w, http.StatusUnauthorized)
		return
	}

	var event Event
	errDecode := json.NewDecoder(r.Body).Decode(&event)
	if errDecode != nil {
		responWithoutBody(w, http.StatusBadRequest)
		return
	}

	if event.Event != "user.upgraded" {
		responWithoutBody(w, http.StatusNoContent)
		return
	}

	errUpdate := apiConfig.dbQueries.UpdateUserStatus(r.Context(), database.UpdateUserStatusParams{
		ID:          event.Data.User_id,
		IsChirpyRed: true,
	})
	if errUpdate != nil {
		responWithoutBody(w, http.StatusNotFound)
		return
	}
	responWithoutBody(w, http.StatusNoContent)

}
