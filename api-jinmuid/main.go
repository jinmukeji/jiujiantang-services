package main

import (
	"os"

	"github.com/jinmukeji/jiujiantang-services/api-jinmuid/config"
	"github.com/jinmukeji/jiujiantang-services/api-jinmuid/rest"
	"github.com/micro/cli"
	"github.com/micro/go-micro/v2/web"
)

var (
	apiBase      string
	jwtSignInKey string
	debug        bool
)

func main() {
	service := web.NewService(
		// Service Basic Info
		web.Name(config.FullServiceName()),
		web.Version(config.ProductVersion),

		// Fault Tolerance - Heartbeating
		web.RegisterTTL(config.DefaultRegisterTTL),
		web.RegisterInterval(config.DefaultRegisterInterval),
		web.Flags(
			cli.StringFlag{
				Name:        "x_api_base",
				Value:       "",
				Usage:       "API Base URL",
				EnvVar:      "X_API_BASE",
				Destination: &apiBase,
			},
			cli.StringFlag{
				Name:        "x_jwt_sign_in_key",
				Usage:       "JWT Sign-in key",
				EnvVar:      "X_JWT_SIGN_IN_KEY",
				Destination: &jwtSignInKey,
			},
			cli.BoolFlag{
				Name:  "version",
				Usage: "Show version information",
			},
			cli.BoolFlag{
				Name:        "x_enable_debug",
				Usage:       "Enable debug",
				EnvVar:      "X_ENABLE_DEBUG",
				Destination: &debug,
			},
		),
	)
	// Init Micro service
	err := service.Init(
		web.Action(func(c *cli.Context) {
			// Setup handler
			app := rest.NewApp(apiBase, jwtSignInKey, debug)
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
