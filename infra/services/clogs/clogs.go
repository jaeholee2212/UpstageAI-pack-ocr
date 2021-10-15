package main

import (
	"fmt"
	"log"
	"net/http"
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

			http.HandleFunc("/clogs/log", clog)

			addr := fmt.Sprintf(":%v", port)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logger.Error(err, "ListenAndServe")
				return err
			}

			return nil
		},
	}
)

func clog(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		logger.Info("not supported", "method", req.Method)
		return
	}

	if ctype, ok := req.Header["Content-Type"]; !ok || len(ctype) <= 0 || ctype[0] != "application/json" {
		logger.Info("Content-Type should be an application/json")
		return
	}

	fmt.Println("req", *req)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		logger.Error(err, "errors from app")
	}
}
