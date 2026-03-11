package api

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	logger     *log.Logger
	httpServer *http.Server
}

func NewServer(logger *log.Logger) *Server {
	r := RouterInit()

	logger.Println("Server initialized")

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		logger:     logger,
		httpServer: httpServer,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
