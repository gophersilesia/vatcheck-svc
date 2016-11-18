package main

import (
	vatcheck "github.com/gopherskatowice/vatcheck-svc"
	. "github.com/gopherskatowice/vatcheck-svc/config"
	"github.com/gopherskatowice/vatcheck-svc/server"
	"github.com/mattes/vat"
	"github.com/urfave/cli"
)

// ServerCommand starts the HTTP server.
func ServerCommand() cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Start the HTTP server",
		Action: func(c *cli.Context) error {
			InitializeConfig()
			InitializeLogging(Config.LogLevel)

			return server.New(
				vatcheck.New(vat.CheckVAT),
				Config.HTTPBind,
				Config.HTTPPort,
			).ListenAndServe()
		},
	}
}
