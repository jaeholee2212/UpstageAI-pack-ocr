package main

import (
	"log"
	"os"

	"github.com/go-logr/stdr"
	"github.com/urfave/cli/v2"
)

var (
	logger = stdr.New(log.New(os.Stdout, "", log.Lshortfile))
	app    = &cli.App{
		Name:  "clogs",
		Usage: "log messages for a snorkel consumption",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "port",
				Usage:    "a port `number`",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			port := c.Int("port")
			logger.Info("clogs", "port", port)
			return nil
		},
	}
)

func main() {
	if err := app.Run(os.Args); err != nil {
		logger.Error(err, "errors from app")
	}
}
