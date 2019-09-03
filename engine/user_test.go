package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type trUserLogin struct {
	ExpCode  int
	Login    string
	Password string
}

func TestLogin(t *testing.T) {

	testRounds := []trUserLogin{
		{
			// empty request
			ExpCode: 400,
		},
		{
			ExpCode:  200,
			Login:    validLogin,
			Password: validPassword,
		},
		{
			ExpCode:  403,
			Login:    validLogin,
			Password: "invalidPassword",
		},
		{
			ExpCode:  403,
			Login:    "invalidLogin",
			Password: validPassword,
		},
		{
			ExpCode:  403,
			Login:    "invalidLogin",
			Password: "invalidPassword",
		},
	}

	env := NewTestEnv(t, "*")
	router := mux.NewRouter()
	env.ExportUserRoutes(router)

	for round, creds := range testRounds {
		jsonStr, err := json.Marshal(creds)
		ok(t, err)
		r, err := http.NewRequest("POST", "/login/", bytes.NewBuffer(jsonStr))
		ok(t, err)

		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)

		assert(t, creds.ExpCode == w.Code, fmt.Sprintf("wrong code `%d` at round %d", w.Code, round))
	}

}
