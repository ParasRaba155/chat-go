// middleware package contains the middleware to be used in this application
package middleware

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"

	response "app/http"
	"app/jwt"
	"app/user"
)

// AuthJWT extracts the token using jwt service and extracts email from the bearer token
// and then stores the user in request context which can letter be used by any controller
func AuthJWT(l *zap.Logger, service jwt.Service, userService user.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearerHeader := r.Header.Get("Authorization")

			if bearerHeader == "" {
				l.Info("empty bearer Header")
				response.SendResponse(w, "No bearer token provided", nil, http.StatusUnauthorized, false)
				return
			}

			headerParts := strings.Split(bearerHeader, " ")
			if len(headerParts) != 2 {
				l.Error("headerparts length is not 2", zap.Int("headerPartsLenght", len(headerParts)))
				response.SendResponse(w, "headerparts length is not 2", nil, http.StatusUnauthorized, false)
				return
			}

			tokenStr := headerParts[1]
			token, err := service.ExtractJWTClaimsFromToken(tokenStr)
			if err != nil {
				l.Error("error in extracting jwt claims", zap.Error(err))
				response.SendResponse(w, "invalid jwt", nil, http.StatusUnauthorized, false)
				return
			}
			email, ok := token["email"].(string)
			if !ok {
				l.Error("error in extracting email", zap.Error(err))
				response.SendResponse(w, "invalid jwt", nil, http.StatusUnauthorized, false)
				return
			}

			user, err := userService.GetUserByEmail(email)
			if err != nil {
				l.Error("error in getting user", zap.Error(err))
				response.SendResponse(w, "invalid jwt", nil, http.StatusUnauthorized, false)
				return
			}

			userCtxKey := struct{}{}
			userCtx := context.WithValue(r.Context(), userCtxKey, user)
			r = r.WithContext(userCtx)
			next.ServeHTTP(w, r)
		})
	}
}

// EnableCors the Cors middleware
func EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // change this later
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}
