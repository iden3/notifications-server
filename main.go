package main

import (
	"os"

	"github.com/iden3/notifications-server/commands"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "notifications-server"
	app.Version = "0.1.0-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
	}

	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, commands.ServerCommands...)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
