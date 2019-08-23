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
