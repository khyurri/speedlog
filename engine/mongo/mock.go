package mongo

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var DelEventsCalledTimes int
var AllEventsCalledTimes int

type DataStoreMock struct {
	FailMetricName         string
	ValidLogin             string
	ValidPassword          string
	DuplicatedProjectTitle string
}

func (d DataStoreMock) FilterEvents(from, to time.Time, metricName, project string) (events []Event, err error) {
	if metricName == d.FailMetricName {
		err = errors.New("[testing] failed by metric name")
		return
	}
	return
}

func (d DataStoreMock) AllEvents(from, to time.Time) (events []AllEvents, err error) {
	if AllEventsCalledTimes == 0 {
		AllEventsCalledTimes = 1
	} else {
		AllEventsCalledTimes++
	}
	return
}

func (d DataStoreMock) SaveEvent(metricName, project string, durationMs float64) (err error) {
	if metricName == d.FailMetricName {
		return errors.New("[testing] failed by metric name")
	}
	return
}

func (d DataStoreMock) AddUser(login string, password string) (err error) {
	panic("implement me")
}

func (d DataStoreMock) GetUser(login string) (*User, error) {
	if login == d.ValidLogin {
		bytes, err := bcrypt.GenerateFromPassword([]byte(d.ValidPassword), 10)
		return &User{
			Login:    d.ValidLogin,
			Password: string(bytes)}, err
	}
	return nil, errors.New("[testing] user not found")
}

func (d DataStoreMock) UserDel(uid string) error {
	panic("implement me")
}

func (d DataStoreMock) AddProject(title string) (err error) {
	if title == d.DuplicatedProjectTitle {
		return errors.New("[testing] failed by project title")
	}
	return
}

func (d DataStoreMock) GetProject(title string) (project Project, err error) {
	panic("implement me")
}

func (d DataStoreMock) GetProjectById(id string) (project Project, err error) {
	panic("implement me")
}

func (d DataStoreMock) DelProject(id string) (err error) {
	panic("implement me")
}

func (d DataStoreMock) DelEvents(to time.Time) (err error) {
	if DelEventsCalledTimes == 0 {
		DelEventsCalledTimes = 1
	} else {
		DelEventsCalledTimes++
	}
	return nil
}
