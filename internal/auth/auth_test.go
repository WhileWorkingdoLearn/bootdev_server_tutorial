package auth_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/WhileCodingDoLearn/bootdev_server_tut/internal/auth"
	"github.com/google/uuid"
)

const secret = "hello world"

func TestJWTToken(t *testing.T) {

	u, _ := uuid.NewUUID()
	token, errToken := auth.MakeJWT(u, secret, 1*time.Second)
	if errToken != nil {
		t.Fatal(errToken)
	}

	u2, errValid := auth.ValidateJWT(token, secret)
	if errValid != nil {
		t.Fatal(errValid)
	}

	if u != u2 {
		t.Failed()
	}

}

func TestGetBeaterToken(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "Bearer adwadwa")
	answer, err := auth.GetToken(header, "Bearer")
	fmt.Println(answer)
	if err != nil {
		t.Fatal(err)
	}
	if answer != "adwadwa" {
		t.Fail()
	}
}

func TestRefreshToken(t *testing.T) {
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		t.Fatal(err)
	}
	if len(refreshToken) == 0 {
		t.Fail()
	}
	fmt.Println(refreshToken)
}
