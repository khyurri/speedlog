package engine

import (
	"github.com/khyurri/speedlog/engine/mongo"
	"log"
)

type Engine struct {
	DBEngine *mongo.Engine
	Logger   *log.Logger
}

func New(dbEngine *mongo.Engine, logger *log.Logger) *Engine {
	return &Engine{dbEngine, logger}
}
