package rest

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
)

type Resp struct {
	Status   int
	JsonBody []byte
	Logger   log.Logger
}

func (r *Resp) Render(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Status {
	case StatusOk:
		w.WriteHeader(http.StatusOK)
	case StatusErr:
		w.WriteHeader(http.StatusBadRequest)
	case StatusIntErr:
		w.WriteHeader(http.StatusInternalServerError)
	case StatusForbidden:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusNoContent)
	}
	_, err := w.Write(r.JsonBody)
	if err != nil {
		r.Logger.Fatal(err)
		return
	}
}
