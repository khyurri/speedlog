package engine

//import (
//	"encoding/json"
//	"net/http"
//)
//
//const (
//	collection = "project"
//)
//
//type RegisterProjectRequest struct {
//	Title string `json:"title"`
//}
//
//func (env *Env) RegisterProjectHttp(w http.ResponseWriter, r *http.Request) {
//	var err error
//
//	response := &Resp{}
//	defer response.Render(w)
//
//	p := &RegisterProjectRequest{}
//	decoder := json.NewDecoder(r.Body)
//	err = decoder.Decode(p)
//	if err != nil {
//		response.Status = StatusErr
//		env.Logger.Fatal(err)
//		return
//	}
//	if len(p.Title) == 0 {
//		response.Status = StatusErr
//		return
//	}
//	_, err = ProjectExists(p.Title, env)
//	// TODO: simplify this block
//	if err != nil {
//		if err.Error() == "not found" {
//			registered := RegisterProject(p.Title, env)
//			if !registered {
//				response.Status = StatusIntErr
//			} else {
//				response.Status = StatusOk
//			}
//		} else {
//			response.Status = StatusIntErr
//		}
//	} else {
//		response.Status = StatusExists
//	}
//
//}
