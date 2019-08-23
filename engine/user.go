package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func invalidAuthRequest(response *Resp) {
	authErrorMsg := "invalid login or password"
	response.Status = StatusForbidden
	response.JsonBody, _ = json.Marshal(
		InvalidRequestParams(authErrorMsg))
}

func tokenResp(token string, response *Resp) {
	response.Status = StatusOk
	response.JsonBody, _ = json.Marshal(struct {
		Token string `json:"token"`
	}{token},
	)
}

func (env *Env) AddUser(login, password string) (err error) {
	err = env.DBEngine.AddUser(login, password)
	if err != nil {
		Logger.Printf("[error] cannot create user %s with password %s", login, password)
	}
	return
}

func (env *Env) authenticateHttp() http.HandlerFunc {

	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		response := &Resp{}
		defer response.Render(w)

		if r.Body == nil {
			_r, _ := json.Marshal(r)
			fmt.Println(_r)
			Logger.Printf("[info] request body is nil. Request: %s", _r)
			invalidAuthRequest(response)
			return
		}
		decoder := json.NewDecoder(r.Body)
		u := &request{}
		err := decoder.Decode(&u)
		if err != nil {
			invalidAuthRequest(response)
			return
		}
		err = env.Authenticate(u.Login, u.Password)
		if err != nil {
			invalidAuthRequest(response)
			return
		}

		// GENERATE TOKEN //
		_, tokenString, err := env.SigningKey.Encode(
			jwt.MapClaims{"source": "rest", "issuer": "hello"})

		if err != nil {
			invalidAuthRequest(response)
		}
		tokenResp(tokenString, response)
	}
}

// Authenticate - returns error, if user not exists or wrong password, else ok
func (env *Env) Authenticate(login string, password string) (err error) {

	u, err := env.DBEngine.GetUser(login)
	if err != nil {
		if err.Error() == "not found" {
			return errors.New("user does not exists")
		}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errors.New("wrong password")
	}
	return
}
