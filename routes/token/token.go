package token

import (
	"net/http"
)

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader != "" {
			// TODO: check token validation
			next(w, req)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}