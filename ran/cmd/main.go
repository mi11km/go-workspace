package main

import (
	"os"

	"github.com/mi11km/go-workspace/ran"
)

func main() {
	app := ran.NewApp()
	os.Exit(app.Run(os.Args, os.Stdin, os.Stdout, os.Stderr))
}
