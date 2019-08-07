package users

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/engine/rest"
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

func AuthenticateHttp(w http.ResponseWriter, r *http.Request, eng *engine.Engine) {
	eng.Logger.Println("Trying to log in")
	decoder := json.NewDecoder(r.Body)
	u := &User{}
	err := decoder.Decode(&u)
	if err != nil {
		eng.Logger.Fatal(err)
		return
	}
	err = Authenticate(u.Login, u.Password, eng)
	if err != nil {
		// TODO: make error response
		eng.Logger.Println(err)
		return
	}
	// GENERATE TOKEN //

	_, tokenString, err := engine.SigningKey.Encode(
		jwt.MapClaims{"source": "rest", "issuer": u.Login})

	if err != nil {
		eng.Logger.Fatal(err)
	}

	eng.Logger.Println(tokenString)
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

func AuthorizeMiddleware(w http.ResponseWriter, r *http.Request) {

}

func Authorize(token string) (result bool) {
	return
}
