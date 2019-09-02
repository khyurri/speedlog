package mongo

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
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
	u := &User{Login: login, Password: hashedPassword}
	if len(u.Login) == 0 || len(u.Password) == 0 {
		err = errors.New("login or password cannot be empty")
		return
	}
	err = mg.Collection(userCollection, sess).Insert(u)
	return
}

// GetUser - returns User by login
func (mg *Mongo) GetUser(login string) (*User, error) {

	sess := mg.Clone()
	defer sess.Close()
	var u *User

	err := mg.Collection(userCollection, sess).Find(bson.M{
		"login": login,
	}).One(&u)
	return u, err
}

// UserDel - delete User by uid (mongodb field _id)
func (mg *Mongo) UserDel(uid string) error {

	sess := mg.Clone()
	defer sess.Close()

	err := mg.Collection(userCollection, sess).Remove(bson.M{
		"_id": bson.ObjectIdHex(uid),
	})

	return err
}
