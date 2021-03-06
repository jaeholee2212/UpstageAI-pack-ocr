package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-logr/stdr"
	"github.com/urfave/cli/v2"
)

var (
	snorkel = NewSnorkel("infra@clogs" /* table name */, "__snorkel-relay__" /* token */).
		AddIntField("elapsed_micros").
		AddIntField("req_bytes").
		AddIntField("http_status").
		AddStrField("method").
		AddStrField("content_type")
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
	stime := time.Now().UnixMicro()

	if req.Method != "POST" {
		logger.Info("not supported", "method", req.Method)
		http.Error(w, "not allowed method", http.StatusMethodNotAllowed)

		snorkel.Write(os.Stdout, EventData{
			"elapsed_micros": time.Now().UnixMicro() - stime,
			"method":         req.Method,
			"req_bytes":      0,
			"http_status":    http.StatusMethodNotAllowed,
		})
		return
	}

	ctype, ok := req.Header["Content-Type"]
	if !ok || len(ctype) <= 0 || ctype[0] != "application/json" {
		logger.Info("Content-Type should be an application/json")
		http.Error(w, "only accept an application/json", http.StatusBadRequest)

		snorkel.Write(os.Stdout, EventData{
			"elapsed_micros": time.Now().UnixMicro() - stime,
			"method":         req.Method,
			"req_bytes":      0,
			"http_status":    http.StatusBadRequest,
			"content_type":   strings.Join(ctype, ","),
		})
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error(err, "failed to get a body of request")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		snorkel.Write(os.Stdout, EventData{
			"elapsed_micros": time.Now().UnixMicro() - stime,
			"method":         req.Method,
			"req_bytes":      0,
			"http_status":    http.StatusInternalServerError,
			"content_type":   strings.Join(ctype, ","),
		})
		return
	}

	// Spit out messages into the stdout so container's log stream receives
	// the given body string which in turn flies into logstash services.
	fmt.Fprintf(os.Stdout, "%v\n", string(body))

	if err := snorkel.Write(os.Stdout, EventData{
		"elapsed_micros": time.Now().UnixMicro() - stime,
		"method":         req.Method,
		"req_bytes":      len(body),
		"http_status":    http.StatusOK,
		"content_type":   strings.Join(ctype, ","),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "snorkel error: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, `{"ok":true}`)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		logger.Error(err, "errors from app")
	}
}
