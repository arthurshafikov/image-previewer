package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Handler interface {
	Init(*gin.Engine)
}

type Logger interface {
	Error(err error)
	Info(msg string)
}

type Server struct {
	logger  Logger
	httpSrv *http.Server
	handler Handler
	Engine  *gin.Engine
}

func NewServer(logger Logger, handler Handler) *Server {
	return &Server{
		logger:  logger,
		handler: handler,
		Engine:  gin.Default(),
	}
}

func (s *Server) Serve(ctx context.Context, g *errgroup.Group, port string) {
	s.handler.Init(s.Engine)

	s.httpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: s.Engine,
	}

	g.Go(func() error {
		<-ctx.Done()
		return s.shutdown()
	})

	if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error(fmt.Errorf("could not start listener: %w", err))
	}
}

func (s *Server) shutdown() error {
	s.logger.Info("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpSrv.Shutdown(ctx)
}
