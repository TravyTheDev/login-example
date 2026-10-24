package users

import (
	"encoding/json"
	"login-example/libs/jwt"
	"net/http"
)

type UserHandler struct {
	store    *UserStore
	jwtMaker *jwt.JWTMaker
}

func NewUserHandler(store *UserStore, jwtMaker *jwt.JWTMaker) *UserHandler {
	return &UserHandler{
		store:    store,
		jwtMaker: jwtMaker,
	}
}

func (h *UserHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /users/me", jwt.GetAuthMiddlewareFunc(h.jwtMaker, h.getUserByID))
	router.HandleFunc("GET /users/all_users", jwt.GetAuthMiddlewareFunc(h.jwtMaker, h.getAllUsers))
}

func (h *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value(jwt.AuthKey{}).(*jwt.UserClaims)

	userID := userClaims.ID

	user, err := h.store.getUserByID(userID)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}

	user.ImpersonatedBy = (*ImpersonatedByUser)(userClaims.ImpersonatedBy)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.getAllUsers()
	if err != nil {
		http.Error(w, "database error loading slips", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "error loading slips", http.StatusInternalServerError)
		return
	}
}
