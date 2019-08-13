package users

import (
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var JWTTestKey = "test_key"

func TestAuthenticateHttp(t *testing.T) {
	logger := log.New(os.Stdout, "speedlog ", log.LstdFlags|log.Lshortfile)
	dbEngine, _ := mongo.New("speedlog", "127.0.0.1:27017", logger)
	eng := engine.New(dbEngine, logger, JWTTestKey)
	errMsg := "{\"message\":\"invalid login or password\"}"

	req, _ := http.NewRequest("POST", "/login/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthenticateHttp(w, r, eng)
	})
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Body.String(), errMsg)
}
