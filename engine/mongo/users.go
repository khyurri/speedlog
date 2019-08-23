package mongo

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Login    string
	Password string
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// AddUser - adds row to database
func (mg *Mongo) AddUser(login string, password string) (err error) {

	sess := mg.Clone()
	defer sess.Close()

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return
	}
	u := &user{Login: login, Password: hashedPassword}
	if len(u.Login) == 0 || len(u.Password) == 0 {
		err = errors.New("login or password cannot be empty")
		return
	}
	err = mg.Collection(userCollection, sess).Insert(u)
	return
}

// GetUser - returns user by login
func (mg *Mongo) GetUser(login string) (*user, error) {

	sess := mg.Clone()
	defer sess.Close()
	var u *user

	err := mg.Collection(userCollection, sess).Find(bson.M{
		"login": login,
	}).One(&u)
	return u, err
}

// UserDel - delete user by uid (mongodb field _id)
func (mg *Mongo) UserDel(uid string) error {

	sess := mg.Clone()
	defer sess.Close()

	err := mg.Collection(userCollection, sess).Remove(bson.M{
		"_id": bson.ObjectIdHex(uid),
	})

	return err
}
