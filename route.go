package main

import (
	"app/chat"
	"app/config"
	"app/jwt"
	"app/middleware"
	"app/user"
	"net/http"

	userQueries "app/user/queries"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Route(r *mux.Router, db *pgxpool.Pool, l *zap.Logger, c *config.Config, rc *redis.Client) {
	jwtService := jwt.NewService(c)
	userRepo := userQueries.New(db)
	userService := user.NewService(userRepo, l)
	redisRepo := user.NewSessionRepo(rc)
	userController := user.NewController(l, &userService, &jwtService, &redisRepo, c.Session)

	// unprotected routes
	r.Use(middleware.EnableCors)
	r.HandleFunc("/api/v1/signup", userController.RegisterUser).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/login", userController.Login).Methods(http.MethodPost)

	// protected routes
	r = r.NewRoute().Subrouter()
	r.Use(middleware.AuthJWT(l, jwtService, userService))

	// user routes
	r.HandleFunc("/api/v1/user/{email}", userController.GetUserByEmail).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/user/reset-password", userController.ResetPassword).Methods(http.MethodPatch)

	room := chat.NewRoom(uuid.New())
	go room.Run(l)
	r.HandleFunc("/api/v1/room/{receiver}", func(w http.ResponseWriter, r *http.Request) {
		chat.Handler(room, l, rc, w, r)
	})
}
