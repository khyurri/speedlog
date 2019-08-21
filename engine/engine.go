package engine

import (
	"github.com/go-chi/jwtauth"
	"github.com/khyurri/speedlog/engine/mongo"
	"log"
)

// Engine - core struct for storing dependencies
type Engine struct {
	DBEngine   *mongo.Engine
	Logger     *log.Logger
	SigningKey *jwtauth.JWTAuth
}

// New - create new engine struct
func New(dbEngine *mongo.Engine, logger *log.Logger, signingKey string) *Engine {
	k := jwtauth.New("HS256", []byte(signingKey), nil)
	return &Engine{dbEngine, logger, k}
}
