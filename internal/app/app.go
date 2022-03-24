package app

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/thewolf27/image-previewer/internal/config"
	"github.com/thewolf27/image-previewer/internal/services"
	server "github.com/thewolf27/image-previewer/internal/transport/http"
	"github.com/thewolf27/image-previewer/internal/transport/http/handler"
)

var (
	configFolder  string
	storageFolder string
)

func init() {
	flag.StringVar(&configFolder, "configFolder", "./configs", "Path to the config folder")
	flag.StringVar(&storageFolder, "storageFolder", "./storage", "Path to the storage folder")
}

func Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := config.NewConfig(configFolder, storageFolder)

	services := services.NewServices(services.Deps{
		Config: config,
	})
	handler := handler.NewHandler(ctx, services)
	server := server.NewServer(handler)

	server.Serve(ctx, config.ServerConfig.Port)
}
