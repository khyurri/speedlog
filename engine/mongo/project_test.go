package mongo

import "testing"

const (
	testProject = "testproject"
)

func TestCRUDProject(t *testing.T) {
	mongo, err := New(testMongoDb, testMongoHost)
	ok(t, err)
	defer mongo.Session.Close()
	defer ok(t, clearDb(t, mongo))

	// Create
	err = mongo.AddProject(testProject)
	ok(t, err)

	// Read
	project, err := mongo.GetProject(testProject)
	ok(t, err)

	project, err = mongo.GetProjectById(project.ID.Hex())
	equals(t, testProject, project.Title)

	// Update is not supported yet

	// Delete
	projectId := project.ID.Hex()
	err = mongo.DelProject(project.ID.Hex())
	project, err = mongo.GetProject(testProject)
	assert(t, err != nil, "project exists!")

	project, err = mongo.GetProjectById(projectId)
	assert(t, err != nil, "project exists!")
}
