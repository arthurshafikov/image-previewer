package app

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/thewolf27/image-previewer/internal/config"
	server "github.com/thewolf27/image-previewer/internal/transport/http"
	"github.com/thewolf27/image-previewer/internal/transport/http/handler"
)

var (
	configFolder string
)

func init() {
	flag.StringVar(&configFolder, "configFolder", "./configs", "Path to the config folder")
}

func Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := config.NewConfig(configFolder)

	handler := handler.NewHandler(ctx)
	server := server.NewServer(handler)

	server.Serve(ctx, config.ServerConfig.Port)
}
