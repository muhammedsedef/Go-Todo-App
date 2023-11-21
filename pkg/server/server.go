package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type server struct {
	engine *gin.Engine
}

func NewServer(engine *gin.Engine) *server {
	return &server{engine: engine}
}

func (s *server) StartHttpServer() {
	server := &http.Server{
		Handler: s.engine,
		Addr:    ":" + "8080",
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("cannot start server", err)
		panic("cannot start server")
	}

	fmt.Println("Server is running on port: 8080")

}
