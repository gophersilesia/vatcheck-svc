package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	. "github.com/gopherskatowice/vatcheck-svc/config"
	"github.com/urfave/cli"
)

var (
	// App codegangsta/cli root cmd
	App *cli.App
)

// Initialize commandline app.
func init() {
	App = cli.NewApp()

	App.Name = "vatcheck-svc"
	App.Usage = `Check if the given VAT number is valid against the EU VIES service`
	App.Author = "Gophers Katowice"

	// Version is injected at build-time
	App.Version = ""

	InitializeConfig()
	InitializeLogging(Config.LogLevel)
}

func main() {
	AddCommands()
	if err := App.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// AddCommands adds child commands to the root command Cmd.
func AddCommands() {
	AddCommand(VersionCommand())
	AddCommand(ServerCommand())
}

// AddCommand adds a child command.
func AddCommand(cmd cli.Command) {
	App.Commands = append(App.Commands, cmd)
}

// InitializeLogging sets logrus log level.
func InitializeLogging(logLevel string) {
	// If log level cannot be resolved, exit gracefully
	if logLevel == "" {
		log.Warning("Log level could not be resolved, fallback to fatal")
		log.SetLevel(log.FatalLevel)
		return
	}
	// Parse level from string
	lvl, err := log.ParseLevel(logLevel)

	if err != nil {
		log.WithFields(log.Fields{
			"passed":  logLevel,
			"default": "fatal",
		}).Warn("Log level is not valid, fallback to default level")
		log.SetLevel(log.FatalLevel)
		return
	}

	log.SetLevel(lvl)
	log.WithFields(log.Fields{
		"level": logLevel,
	}).Debug("Log level successfully set")
}
