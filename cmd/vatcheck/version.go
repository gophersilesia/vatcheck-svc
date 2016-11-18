package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli"
)

// Do not set these manually! these variables
// are meant to be set through ldflags.
var buildTag, buildDate string

// VersionCommand returns a version command.
func VersionCommand() cli.Command {
	return cli.Command{
		Name:  "version",
		Usage: "Print the version number of the service",
		Action: func(ctx *cli.Context) error {
			if buildTag != "" && buildDate != "" {
				fmt.Printf("%s built on %s\n", buildTag, buildDate)
				return nil
			}
			return errors.New("No version")
		},
	}
}
