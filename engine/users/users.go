package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/rest"
	"net/http"
)

const (
	collection = "users"
)

type User struct {
	Login    string
	Password string
}

func ExportRoutes(router *mux.Router, app *rest.App) {
	router.HandleFunc("/login/", app.MongoEngine(AuthenticateHttp)).
		Methods("POST")
}

func invalidAuthRequest(response *rest.Resp) {
	authErrorMsg := "invalid login or password"
	response.Status = rest.StatusErr
	response.JsonBody, _ = json.Marshal(
		rest.InvalidRequestParams(authErrorMsg))
}

func tokenResp(token string, response *rest.Resp) {
	response.Status = rest.StatusOk
	response.JsonBody, _ = json.Marshal(struct {
		Token string `json:"token"`
	}{token},
	)
}

func AuthenticateHttp(w http.ResponseWriter, r *http.Request, eng *engine.Engine) {
	response := &rest.Resp{}
	defer response.Render(w)

	if r.Body == nil {
		_r, _ := json.Marshal(r)
		fmt.Println(_r)
		eng.Logger.Printf("[info] request body is nil. Request: %s", _r)
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
	err = Authenticate(u.Login, u.Password, eng)
	if err != nil {
		invalidAuthRequest(response)
		return
	}

	// GENERATE TOKEN //
	eng.Logger.Println("----------------------------")
	eng.Logger.Println(u.Login)
	eng.Logger.Println("----------------------------")
	_, tokenString, err := eng.SigningKey.Encode(
		jwt.MapClaims{"source": "rest", "issuer": u.Login})

	if err != nil {
		invalidAuthRequest(response)
	}
	tokenResp(tokenString, response)
}

func AddUser(login string, password string, eng *engine.Engine) (err error) {
	u := &User{login, password}
	if len(u.Login) == 0 || len(u.Password) == 0 {
		err = errors.New("login or password cannot be empty")
		return
	}
	err = eng.DBEngine.Collection(collection).Insert(u)
	return
}

func Authenticate(login string, password string, eng *engine.Engine) (err error) {
	var u User
	err = eng.DBEngine.Collection("users").Find(bson.M{
		"login":    login,
		"password": password,
	}).One(&u)
	if err != nil {
		if err.Error() == "not found" {
			return errors.New("user does not exists")
		} else {
			eng.Logger.Fatal(err)
			return
		}
	}
	return
}
