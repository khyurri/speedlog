package mongo

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Login    string
	Password string
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func isMatch(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// AddUser - adds row to database
func (mg *Mongo) AddUser(login string, password string) (err error) {

	sess := mg.Clone()
	defer sess.Close()

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return
	}
	u := &user{login, hashedPassword}
	if len(u.Login) == 0 || len(u.Password) == 0 {
		err = errors.New("login or password cannot be empty")
		return
	}
	err = mg.Collection(userCollection, sess).Insert(u)
	return
}

// Authenticate - returns error, if user not exists or wrong password
func (mg *Mongo) Authenticate(login string, password string) (err error) {
	sess := mg.Clone()
	defer sess.Close()

	var u user
	err = mg.Collection("users", sess).Find(bson.M{
		"login": login,
	}).One(&u)
	if err != nil {
		if err.Error() == "not found" {
			return errors.New("user does not exists")
		}
		return
	}
	if !isMatch(u.Password, password) {
		return errors.New("wrong password")
	}
	return
}
