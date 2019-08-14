package events

import (
	"github.com/khyurri/speedlog/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EventsTestSuit struct {
	suite.Suite
	test.SpeedLogTest
	TestEvents []interface{}
}

func (suite *EventsTestSuit) SetupTest() {
	suite.Init()
}

func TestEvents(t *testing.T) {
	suite.Run(t, new(EventsTestSuit))
}
