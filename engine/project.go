package engine

import (
	"encoding/json"
	"fmt"
	"github.com/khyurri/speedlog/utils"
	"net/http"
)

func (env *Env) createProjectHttp() http.HandlerFunc {

	type request struct {
		Title string `json:"title"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		response := &Resp{}
		defer response.Render(w)

		req := &request{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(req)
		utils.Debug(fmt.Sprintf("trying to create project"))
		if len(req.Title) == 0 {
			response.Status = StatusErr
			return
		}
		err = env.DBEngine.AddProject(req.Title)
		if err != nil {
			// TODO: check if dublicate
			response.Status = StatusErr
			utils.Debug(fmt.Sprintf("create project %s", err))
		}
	}
}
