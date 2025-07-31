package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	addr   string
	engine *gin.Engine
}

func New(addr string) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	return &Server{
		addr:   addr,
		engine: engine,
	}
}

func (s *Server) SetupRoutes() {
	s.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	s.engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/metrics")
	})
}

func (s *Server) Start() error {
	return s.engine.Run(s.addr)
}
