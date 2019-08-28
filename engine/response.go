package engine

import (
	"log"
	"net/http"
)

// http functions status
const (
	StatusOk        = 0
	StatusIntErr    = 1
	StatusErr       = 2
	StatusForbidden = 3
	StatusExists    = 4
)

// struct for http response
type Resp struct {
	Status   int
	JsonBody []byte
	Logger   log.Logger
}

func (r *Resp) setHeader(w http.ResponseWriter) {
	switch r.Status {
	case StatusOk:
		w.WriteHeader(http.StatusOK)
	case StatusErr:
		w.WriteHeader(http.StatusBadRequest)
	case StatusIntErr:
		w.WriteHeader(http.StatusInternalServerError)
	case StatusForbidden:
		w.WriteHeader(http.StatusForbidden)
	case StatusExists:
		w.WriteHeader(http.StatusNotModified)
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

// just render JSON response from struct Resp
func (r *Resp) Render(w http.ResponseWriter) {
	r.setHeader(w)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(r.JsonBody)
	if err != nil {
		// TODO: fix null pointer exception
		// r.Logger.Fatal(err)
		return
	}
}

// returns struct ready for rendering with text message
func InvalidRequestParams(message string) interface{} {
	return struct {
		Message string `json:"message"`
	}{message}
}
