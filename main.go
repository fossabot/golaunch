package main

import (
	"os"

	"github.com/micnncim/golaunch/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "golaunch"
	app.Usage = "Launch app"
	app.Action = cmd.Launch

	app.Run(os.Args)
}
