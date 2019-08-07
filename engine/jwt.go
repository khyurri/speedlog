package engine

import (
	"fmt"
	"github.com/go-chi/jwtauth"
	"net/http"
)

var SigningKey = jwtauth.New("HS256", []byte("HelloKey"), nil)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := jwtauth.TokenFromHeader(r)
		f, err := SigningKey.Decode(t)
		fmt.Println(f)
		if err != nil {
			// TODO: 403
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
