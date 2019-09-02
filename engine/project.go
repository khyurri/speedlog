package engine

import (
	"encoding/json"
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
		env.Logger.Println("[debug] trying to create project")
		if len(req.Title) == 0 {
			response.Status = StatusErr
			return
		}
		err = env.DBEngine.AddProject(req.Title)
		if err != nil {
			// TODO: check if dublicate
			response.Status = StatusErr
			env.Logger.Printf("[error] create project %s", err)
		}
	}
}
