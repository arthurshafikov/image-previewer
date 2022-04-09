package app

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/arthurshafikov/image-previewer/internal/config"
	"github.com/arthurshafikov/image-previewer/internal/image_cache"
	"github.com/arthurshafikov/image-previewer/internal/services"
	server "github.com/arthurshafikov/image-previewer/internal/transport/http"
	"github.com/arthurshafikov/image-previewer/internal/transport/http/handler"
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
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := config.NewConfig(configFolder, storageFolder)

	rawImageCache := image_cache.NewCache(config.AppConfig.SizeOfLRUCacheForRawImages, storageFolder+"/raw")
	resizedImageCache := image_cache.NewCache(config.AppConfig.SizeOfLRUCacheForResizedImages, storageFolder+"/resized")
	services := services.NewServices(services.Deps{
		Config:            config,
		RawImageCache:     rawImageCache,
		ResizedImageCache: resizedImageCache,
	})
	handler := handler.NewHandler(ctx, services)
	server := server.NewServer(handler)

	server.Serve(ctx, config.ServerConfig.Port)
}
