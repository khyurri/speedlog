package engine

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func invalidAuthRequest(response *Resp) {
	authErrorMsg := "invalid login or password"
	response.Status = StatusErr
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

func (env *Env) AuthenticateHttp(w http.ResponseWriter, r *http.Request) {
	response := &Resp{}
	defer response.Render(w)

	if r.Body == nil {
		_r, _ := json.Marshal(r)
		fmt.Println(_r)
		env.Logger.Printf("[info] request body is nil. Request: %s", _r)
		invalidAuthRequest(response)
		return
	}
	decoder := json.NewDecoder(r.Body)
	u := &User{}
	err := decoder.Decode(&u)
	if err != nil {
		invalidAuthRequest(response)
		return
	}
	err = env.DBEngine.Authenticate(u.Login, u.Password, env)
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
