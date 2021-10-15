package main

import (
	"log"
	"os"

	"github.com/go-logr/stdr"
)

var (
	logger = stdr.New(log.New(os.Stdout, "", log.Lshortfile))
)

func main() {
	logger.Info("Hello Clogs!")
}
