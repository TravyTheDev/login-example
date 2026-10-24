package auth

import (
	"encoding/json"
	"fmt"
	"login-example/libs/jwt"
	"login-example/utils"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

const accessDuration = time.Minute * 15
const refreshDuration = 24 * time.Hour

type AuthHandler struct {
	store    *AuthStore
	jwtMaker *jwt.JWTMaker
}

func NewAuthHandler(store *AuthStore, jwtMaker *jwt.JWTMaker) *AuthHandler {
	return &AuthHandler{
		store:    store,
		jwtMaker: jwtMaker,
	}
}

func (h *AuthHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /auth/register", h.register)
	router.HandleFunc("POST /auth/login", h.login)
	router.HandleFunc("POST /auth/renew_token", h.renewAccessToken)
	router.HandleFunc("POST /auth/logout", h.logout)
	router.HandleFunc("POST /auth/login_as", jwt.GetAuthMiddlewareFunc(h.jwtMaker, h.loginAs))
	router.HandleFunc("POST /auth/logout_from_loginAs", jwt.GetAuthMiddlewareFunc(h.jwtMaker, h.logutFromloginAs))
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
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
	u := &UserPayload{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	cookie := jwt.GetCookieHandler(w, r, "refresh")
	claims, err := h.jwtMaker.VerifyToken(cookie)
	if err != nil {
		http.Error(w, "error verifying token", http.StatusInternalServerError)
		return
	}

	if err := h.store.DeleteSession(claims.RegisteredClaims.ID); err != nil {
		http.Error(w, "error deleting session", http.StatusInternalServerError)
		return
	}
	jwt.DeleteCookieHandler(w, r, "authentication")
	jwt.DeleteCookieHandler(w, r, "refresh")
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) loginAs(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value(jwt.AuthKey{}).(*jwt.UserClaims)
	var payload LoginAsPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "login as error", http.StatusInternalServerError)
		return
	}

	role := userClaims.Role
	if role != 0 {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	user, err := h.store.GetUserByEmail(payload.TargetEmail)
	if err != nil {
		http.Error(w, "error retrieving target user", http.StatusBadRequest)
		return
	}

	impersonatedByUser := jwt.ImpersonatedByUser{
		ID:       userClaims.ID,
		Username: userClaims.Username,
		Email:    userClaims.Email,
	}

	accessToken, _, err := h.jwtMaker.CreateImpersonationToken(
		int(user.ID),
		user.Username,
		user.Email,
		user.Role,
		impersonatedByUser,
		accessDuration,
	)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	jwt.SetCookieHandler(w, r, accessToken, int(accessDuration.Seconds()), "authentication")
}

func (h *AuthHandler) logutFromloginAs(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value(jwt.AuthKey{}).(*jwt.UserClaims)
	email := userClaims.ImpersonatedBy.Email

	user, err := h.store.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusBadRequest)
		return
	}

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

	jwt.SetCookieHandler(w, r, accessToken, int(accessDuration.Seconds()), "authentication")
	u := &UserPayload{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) renewAccessToken(w http.ResponseWriter, r *http.Request) {

	refreshToken := jwt.GetCookieHandler(w, r, "refresh")

	refreshClaims, err := h.jwtMaker.VerifyToken(refreshToken)
	if err != nil {
		http.Error(w, "error verifying token", http.StatusUnauthorized)
		return
	}

	session, err := h.store.GetSession(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		http.Error(w, "error getting sessions", http.StatusInternalServerError)
		return
	}

	if session.IsRevoked {
		http.Error(w, "session revoked", http.StatusUnauthorized)
		return
	}

	if session.UserEmail != refreshClaims.Email {
		http.Error(w, "session revoked", http.StatusUnauthorized)
		return
	}

	var duration = time.Minute * 15

	accessToken, _, err := h.jwtMaker.CreateToken(
		refreshClaims.ID,
		refreshClaims.Username,
		refreshClaims.Email,
		refreshClaims.Role,
		duration,
	)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	jwt.SetCookieHandler(w, r, accessToken, int(duration.Seconds()), "authentication")

}

func (h *AuthHandler) setCookiesAndSession(user *User, w http.ResponseWriter, r *http.Request) {
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

	jwt.SetCookieHandler(w, r, accessToken, int(accessDuration.Seconds()), "authentication")
	jwt.SetCookieHandler(w, r, refreshToken, int(refreshDuration.Seconds()), "refresh")
}
