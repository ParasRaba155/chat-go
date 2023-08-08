// user domain handles all the user related controllers, services and db logic including session
package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"app/config"
	httputil "app/http"
	"app/jwt"
	"app/user/queries"
)

type Controller struct {
	UserService   *Service
	JwtService    *jwt.Service
	RedisRepo     *sessionRepo
	SessionConfig config.Session
	Logger        *zap.Logger
}

func NewController(l *zap.Logger, u *Service, j *jwt.Service, sr *sessionRepo, cfg config.Session) *Controller {
	return &Controller{
		UserService:   u,
		Logger:        l,
		JwtService:    j,
		RedisRepo:     sr,
		SessionConfig: cfg,
	}
}

func (c *Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req queries.RegisterUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.Logger.Error("could not read request", zap.String("operation", "RegisterUser"), zap.Error(err))
		httputil.SendResponse(w, "could not read the request", nil, http.StatusBadRequest, false)
		return
	}
	err = c.UserService.RegisterUser(req)
	if err != nil {
		c.Logger.Error("could not read operation", zap.String("operation", "RegisterUser"), zap.Error(err))
		httputil.SendResponse(w, err.Error(), nil, http.StatusBadRequest, false)
		return
	}
	httputil.SendResponse(w, "successfully registered the user", nil, http.StatusCreated, true)
	c.Logger.Info("Registered the user")
}

func (c *Controller) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email, ok := mux.Vars(r)["email"]
	if !ok {
		email = ""
	}
	user, err := c.UserService.GetUserByEmail(email)
	if err != nil {
		c.Logger.Error("could not get user", zap.Error(err))
		httputil.SendResponse(w, "user does not exist", nil, http.StatusBadRequest, false)
		return
	}
	httputil.SendResponse(w, "found the user with given email", user, 200, true)
	c.Logger.Info("Found the user", zap.String("email", email))
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.Logger.Error("could not read request", zap.String("operation", "UserLogin"), zap.Error(err))
		httputil.SendResponse(w, "could not read the request", nil, http.StatusBadRequest, false)
		return
	}
	user, err := c.UserService.GetUser(req.Email, req.Password)
	if err != nil {
		c.Logger.Error("could not get User", zap.String("operation", "UserLogin"), zap.Error(err))
		httputil.SendResponse(w, err.Error(), nil, http.StatusBadRequest, false)
		return
	}
	token, err := c.JwtService.GenerateJWT(user.Email)
	if err != nil {
		c.Logger.Error("could not generate jwt token", zap.String("operation", "UserLogin"), zap.Error(err))
		httputil.SendResponse(w, err.Error(), nil, http.StatusBadRequest, false)
		return
	}
	sessionKey, err := c.RedisRepo.CreateSession(r.Context(), req.Email, c.SessionConfig.MaxAgeSec)
	if err != nil {
		c.Logger.Error("could not generate session key", zap.String("operation", "UserLogin"), zap.Error(err))
		httputil.SendResponse(w, err.Error(), nil, http.StatusInternalServerError, false)
		return
	}
	cookie := httputil.CreateCookie(c.SessionConfig.CookieName, sessionKey, c.SessionConfig.MaxAgeSec)
	http.SetCookie(w, cookie)

	resp := user.ToUserWithToken(token)

	httputil.SendResponse(w, "successfully registered the user", resp, http.StatusOK, true)
	c.Logger.Info("Logged Successfully", zap.String("email", resp.Email))
}

func (c *Controller) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CurrentPass string `json:"currentPassword"`
		NewPass     string `json:"newPassword"`
		ConfirmPass string `json:"confirmPassword"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		c.Logger.Error("could not read request", zap.String("operation", "ResetPassword"), zap.Error(err))
		httputil.SendResponse(w, "could not read the request", nil, http.StatusBadRequest, false)
		return
	}
	user := httputil.GetUserFromRequestContext(r)
	err = c.UserService.ResetPassword(user.Email, req.CurrentPass, req.NewPass, req.ConfirmPass)
	if err != nil {
		c.Logger.Error("could not reset password", zap.String("operation", "ResetPassword"), zap.Error(err))
		httputil.SendResponse(w, err.Error(), nil, http.StatusBadRequest, false)
		return
	}

	httputil.SendResponse(w, "you have successfully changed your password", nil, http.StatusAccepted, true)
}
