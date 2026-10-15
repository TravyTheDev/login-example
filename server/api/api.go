package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/cors"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	frontUrl := os.Getenv("FRONT_URL")
	secretKey := os.Getenv("JWT_KEY")
	router := http.NewServeMux()

	userStore := users.NewUserStore(s.db)
	usersHandler := users.NewUserHandler(userStore, secretKey)
	usersHandler.RegisterRoutes(router)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{frontUrl},
		AllowCredentials: true,
		AllowedMethods:   []string{"OPTIONS", http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
	})
	handler := c.Handler(router)

	fmt.Printf("Listening on %s\n", s.addr)
	return http.ListenAndServe(s.addr, handler)
}
