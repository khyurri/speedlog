package engine

import (
	"github.com/go-chi/jwtauth"
	"net/http"
)

func (env *Env) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := jwtauth.TokenFromHeader(r)
		_, err := env.SigningKey.Decode(t)
		if err != nil {
			response := &Resp{}
			response.Status = StatusForbidden
			response.Render(w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (env *Env) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(env.AllowOrigin) > 0 {
			Logger.Printf("[debug] Access-Control-Allow-Origin: %s", env.AllowOrigin)
			w.Header().Set("Access-Control-Allow-Origin", env.AllowOrigin)
		}
		switch r.Method {
		case "OPTIONS":
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			break
		default:
			next.ServeHTTP(w, r)
		}
	})
}
