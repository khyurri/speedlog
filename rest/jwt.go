package rest

//
//import (
//	"github.com/go-chi/jwtauth"
//	"net/http"
//)
//
//var SigningKey = jwtauth.New("HS256", []byte("HelloKey"), nil)
//
//func JWTMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		t := jwtauth.TokenFromHeader(r)
//		_, err := SigningKey.Decode(t)
//		if err != nil {
//			response := &Resp{}
//			response.Status = StatusForbidden
//			response.Render(w)
//		} else {
//			next.ServeHTTP(w, r)
//		}
//	})
//}
