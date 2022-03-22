package app

import (
	"flag"

	"github.com/thewolf27/image-previewer/internal/config"
)

var (
	configFolder string
)

func init() {
	flag.StringVar(&configFolder, "configFolder", "./configs", "Path to the config folder")
}

func Run() {
	_ = config.NewConfig(configFolder)
}
