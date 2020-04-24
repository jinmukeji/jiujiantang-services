package main

import (
	"os"

	"github.com/jinmukeji/jiujiantang-services/web-go/config"
	"github.com/jinmukeji/jiujiantang-services/web-go/rest"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/web"
)

var (
	configFile string
)

func main() {
	service := web.NewService(
		// Service Basic Info
		web.Name(config.FullServiceName()),
		web.Version(config.ProductVersion),

		// Fault Tolerance - Heartbeating
		web.RegisterTTL(config.DefaultRegisterTTL),
		web.RegisterInterval(config.DefaultRegisterInterval),

		webOptions(),
	)
	// Setup --version flag

	// Init Micro service
	err := service.Init(
		web.Action(func(c *cli.Context) {
			// Setup handler
			app := rest.NewApp(configFile)
			service.Handle("/", app)

			if c.Bool("version") {
				config.PrintFullVersionInfo()
				os.Exit(0)
			}
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Service Name:", config.FullServiceName())
	log.Infoln("Version:", config.ProductVersion)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func webOptions() web.Option {
	return web.Flags(
		&cli.StringFlag{
			Name:        "x_config_file",
			Usage:       "Config File",
			EnvVars:     []string{"X_CONFIG_FILE"},
			Destination: &configFile,
		},
		&cli.BoolFlag{
			Name:  "version",
			Usage: "Show version information",
		},
	)
}
