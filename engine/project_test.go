package engine

import (
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

type ProjectTestSuit struct {
	suite.Suite
}

func (suite *ProjectTestSuit) SetupTest() {
	Logger = log.New(os.Stdout, "speedlog ", log.LstdFlags|log.Lshortfile)
}

func (suite *ProjectTestSuit) TestCreateProject() {
	project := "test_project"
	dbEngine, _ := mongo.New("speedlog", "127.0.0.1:27017")
	err := dbEngine.AddProject(project)
	assert.Nil(suite.T(), err)

	projectId, err := dbEngine.GetProject(project)
	assert.Nil(suite.T(), err)
	assert.Greater(suite.T(), len(projectId), 0)

	err = dbEngine.DelProject(projectId)
	assert.Nil(suite.T(), err)

}

func TestProject(t *testing.T) {
	suite.Run(t, new(ProjectTestSuit))
}
