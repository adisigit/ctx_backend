package server

import (
	"ctx_backend/internal/auth"
	"ctx_backend/internal/auth/providers"
	"ctx_backend/internal/database"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	auth.RegisterProvider("google", &providers.GoogleProvider{})
	NewServer := &Server{
		port: port,
		db:   database.New(),
	}
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return server
}
