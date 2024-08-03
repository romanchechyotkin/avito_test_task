package httpsrv

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/romanchechyotkin/avito_test_task/internal/config"

	"github.com/gin-gonic/gin"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	log    *slog.Logger
	cfg    *config.Config
	notify chan error

	router *gin.Engine
	base   *http.Server
}

func New(log *slog.Logger, cfg *config.Config) (*Server, error) {
	router := gin.Default()

	srv := &Server{
		log:    log,
		cfg:    cfg,
		notify: make(chan error),
		router: router,
		base: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
			Handler:      router,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
		},
	}

	srv.start()

	return srv, nil
}

func (srv *Server) start() {
	go func() {
		srv.notify <- srv.base.ListenAndServe()
		close(srv.notify)
	}()
}

func (srv *Server) Notify() <-chan error {
	return srv.notify
}

func (srv *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	return srv.base.Shutdown(ctx)
}

func (srv *Server) RegisterRoutes() {
	srv.router.GET("/status", func(c *gin.Context) {
		c.String(http.StatusOK, "sad")
	})

	authGroup := srv.router.Group("/auth")
	authGroup.POST("/registration", srv.Registration)
	authGroup.POST("/login", srv.Login)
}
