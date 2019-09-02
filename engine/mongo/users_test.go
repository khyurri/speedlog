package mongo

import (
	"fmt"
	"testing"
)

const (
	validLogin    = "vLogin"
	validPassword = "vPassword"
)

func TestCrudUsers(t *testing.T) {
	mongo, err := New(testMongoDb, testMongoHost)
	ok(t, err)
	defer mongo.Session.Close()
	defer ok(t, clearDb(t, mongo))

	// create
	err = mongo.AddUser(validLogin, validPassword)
	ok(t, err)

	err = mongo.AddUser("", "")
	assert(t, err != nil, fmt.Sprintf("login or password cannot be empty"))

	// read
	user, err := mongo.GetUser(validLogin)
	ok(t, err)
	equals(t, validLogin, user.Login)

	// update is not supported yet

	// delete
	err = mongo.UserDel(user.Id.Hex())
	ok(t, err)

}
