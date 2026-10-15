package users

import (
	"encoding/json"
	"fmt"
	"login-example/service/jwt"
	"login-example/utils"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

const accessDuration = time.Minute * 15
const refreshDuration = 24 * time.Hour

type UserHandler struct {
	store    *UserStore
	jwtMaker *jwt.JWTMaker
}

func NewUserHandler(store *UserStore, secretKey string) *UserHandler {
	return &UserHandler{
		store:    store,
		jwtMaker: jwt.NewJWTMaker(secretKey),
	}
}

func (h *UserHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /users/register", h.register)
	router.HandleFunc("POST /users/login", h.login)
}

func (h *UserHandler) register(w http.ResponseWriter, r *http.Request) {
	var payload RegisterPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println(err)
	}

	if payload.Password != payload.PasswordConfirm {
		http.Error(w, "passwords must match", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		http.Error(w, "error handling password", http.StatusInternalServerError)
		return
	}

	user := User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashedPassword,
		Role:     0,
	}

	if err := h.store.CreateUser(user); err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "login error", http.StatusInternalServerError)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		if errors != nil {
			http.Error(w, "need email and password", http.StatusBadRequest)
		}
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusBadRequest)
		return
	}
	h.setCookiesAndSession(user, w, r)
	w.Header().Add("Content-Type", "application/json")
}

func (h *UserHandler) setCookiesAndSession(user *User, w http.ResponseWriter, r *http.Request) {
	accessToken, _, err := h.jwtMaker.CreateToken(
		int(user.ID),
		user.Username,
		user.Email,
		user.Role,
		accessDuration,
	)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}
	refreshToken, refreshClaims, err := h.jwtMaker.CreateToken(
		int(user.ID),
		user.Username,
		user.Email,
		user.Role,
		refreshDuration,
	)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}
	sesh := &Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    user.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    refreshClaims.ExpiresAt.Time,
	}

	_, err = h.store.CreateSession(sesh)
	if err != nil {
		http.Error(w, "error creating session", http.StatusInternalServerError)
		return
	}

	jwt.SetCookieHandler(w, r, accessToken, int(accessDuration), "authentication")
	jwt.SetCookieHandler(w, r, refreshToken, int(refreshDuration), "refresh")
}
