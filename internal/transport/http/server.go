package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Handler interface {
	Init(*gin.Engine)
}

type Server struct {
	httpSrv *http.Server
	handler Handler
	Engine  *gin.Engine
}

func NewServer(handler Handler) *Server {
	return &Server{
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
		log.Println("Could not start listener ", err)
	}
}

func (s *Server) shutdown() error {
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpSrv.Shutdown(ctx)
}
