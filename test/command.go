package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "commander"
	app.Usage = "test command args!"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Value: "sumory",
			Usage: "my name",
		},
		cli.StringFlag{
			Name:  "age",
			Value: "18",
			Usage: "my age",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "set",
			Aliases: []string{"s"},
			Usage:   "set key/value",
			//Flags:   app.Flags,
			Action: func(c *cli.Context) {
				for i, k := range c.GlobalFlagNames() {
					fmt.Printf("%d %s:%v\n", i, k,c.GlobalString(k))
				}
				println("set k/v: ", c.Args().Get(0), c.Args().Get(1))
			},
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "get value of the given key",
			Flags:   app.Flags,
			Action: func(c *cli.Context) {
				fmt.Printf("%+v\n", c.Args())
				println("get key: ", c.Args().Get(0))
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "delete the given key",
			Flags:   app.Flags,
			Action: func(c *cli.Context) {
				fmt.Printf("%+v\n", c.Args())
				println("delete key: ", c.Args().Get(0))
			},
		},
	}

	app.Run(os.Args)
}
