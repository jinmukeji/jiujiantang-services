package main

import (
	"os"

	"github.com/jinmukeji/jiujiantang-services/api-sys/config"
	"github.com/jinmukeji/jiujiantang-services/api-sys/rest"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/web"
)

var (
	apiBase    string
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
			app := rest.NewApp(apiBase, configFile)
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
	log.Infof("API Base: /%s", apiBase)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func webOptions() web.Option {
	return web.Flags(
		&cli.StringFlag{
			Name:        "x_api_base",
			Value:       "",
			Usage:       "API Base URL",
			EnvVars:     []string{"X_API_BASE"},
			Destination: &apiBase,
		},
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
