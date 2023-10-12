package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lusingander/zsx"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:  "zsx",
		Usage: "aws utils",
		Commands: []*cli.Command{
			{
				Name:   "list-profiles",
				Usage:  "`aws configure list-profiles`",
				Action: listProfilesAction,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func listProfilesAction(cCtx *cli.Context) error {
	profiles, err := zsx.ListProfiles()
	if err != nil {
		return err
	}
	for _, p := range profiles {
		fmt.Println(p)
	}
	return nil
}
