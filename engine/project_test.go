package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type trProjCreate struct {
	Title   string `json:"title"`
	Login   string
	ExpCode int
}

func TestCreateProject(t *testing.T) {
	testRounds := []trProjCreate{
		{
			ExpCode: 403,
		},
		{
			Login:   validLogin,
			ExpCode: 400,
		},
		{
			Login:   validLogin,
			Title:   duplicatedProjectTitle,
			ExpCode: 400,
		},
		{
			Login:   validLogin,
			Title:   "validTitle",
			ExpCode: 200,
		},
	}

	env := NewTestEnv(t, "*")
	router := mux.NewRouter()
	env.ExportProjectRoutes(router)

	for round, project := range testRounds {
		jsonStr, err := json.Marshal(project)
		ok(t, err)
		r, err := http.NewRequest("PUT", "/private/project/", bytes.NewBuffer(jsonStr))
		ok(t, err)

		// authorize
		if len(project.Login) > 0 {
			token := getToken(t, env, project.Login)
			r.Header.Set("Authorization", "Bearer "+token)
		}

		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		assert(t, project.ExpCode == w.Code, fmt.Sprintf("wrong code `%d` at round %d", w.Code, round))
		options(t, "/private/project/", router)
	}
}
