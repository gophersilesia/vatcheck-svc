package main

import (
	. "github.com/gopherskatowice/vatcheck-svc/config"
	"github.com/gopherskatowice/vatcheck-svc/server"
	"github.com/urfave/cli"
)

// ServerCommand starts the HTTP server.
func ServerCommand() cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Start the HTTP server",
		Action: func(c *cli.Context) error {
			srv := server.New(
				Config.HTTPBind,
				Config.HTTPPort,
			)
			return srv.ListenAndServe()
		},
	}
}
