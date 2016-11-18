package main

import (
	"github.com/codegangsta/cli"
	. "github.com/gopherskatowice/vatcheck-svc/config"
	"github.com/gopherskatowice/vatcheck-svc/server"
)

// ServerCommand starts the HTTP server.
func ServerCommand() cli.Command {
	return cli.Command{
		Name:    "server",
		Aliases: []string{"serve", "api"},
		Usage:   "Start the HTTP server",
		Action: func(ctx *cli.Context) {
			srv := server.New(
				Config.HTTPBind,
				Config.HTTPPort,
				Config.DisableCors,
			)
			srv.ListenAndServe()
		},
	}
}
