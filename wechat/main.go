package main

import (
	"os"

	"github.com/jinmukeji/jiujiantang-services/wechat/config"
	"github.com/jinmukeji/jiujiantang-services/wechat/rest"
	"github.com/micro/cli"
	"github.com/micro/go-micro/web"
)

var (
	// REST 参数
	options rest.Options
)

func main() {
	service := web.NewService(
		// Service Basic Info
		web.Name(config.FullServiceName()),
		web.Version(config.ProductVersion),

		// Fault Tolerance - Heartbeating
		web.RegisterTTL(config.DefaultRegisterTTL),
		web.RegisterInterval(config.DefaultRegisterInterval),

		// CLI Flags
		// Setup --version flag
		web.Flags(
			cli.StringFlag{
				Name:        "x_api_base",
				Value:       "",
				Usage:       "API Base URL",
				EnvVar:      "X_API_BASE",
				Destination: &(options.APIBase),
			},
			cli.StringFlag{
				Name:        "x_wx_callback_server_base",
				Value:       "",
				Usage:       "微信公众平台 回调服务器的地址",
				EnvVar:      "X_WX_CALLBACK_SERVER_BASE",
				Destination: &(options.WxCallbackServerBase),
			},
			cli.StringFlag{
				Name:        "x_wx_h5_server_base",
				Value:       "",
				Usage:       "微信公众平台 模版ID",
				EnvVar:      "X_WX_H5_SERVER_BASE",
				Destination: &(options.WxH5ServerBase),
			},
			cli.BoolFlag{
				Name:  "version",
				Usage: "Show version information",
			},
		),
	)

	// Init Micro service
	err := service.Init(
		web.Action(func(c *cli.Context) {
			// Setup handler
			app := rest.NewApp(&options)
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
	log.Infof("API Base: /%s", options.APIBase)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
