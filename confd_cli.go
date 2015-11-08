package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	mycli "github.com/sumory/confd/cli"
	"github.com/sumory/confd/config"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     false,
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05.00000",
	})
	log.SetLevel(log.DebugLevel)
}

func parseArgs() {
	app := cli.NewApp()
	app.Name = "confd cli tool"
	app.Usage = "use this tool to operate confd store!"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "store-type",
			Value: "file",
			Usage: "backend store to use",
		},
		cli.StringFlag{
			Name:  "confdir",
			Value: "/data/confd",
			Usage: "confd conf directory",
		},
		cli.StringFlag{
			Name:  "config-file",
			Value: "/data/confd/data/config.toml",
			Usage: "the confd config file",
		},
		cli.StringFlag{
			Name:  "connect-addr",
			Value: "",
			Usage: "backend store address",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "getall",
			Aliases: []string{"gl"},
			Usage:   "get all key/value, e.g. `confd-cli getall`",
			Action: func(c *cli.Context) {
				getall(c)
			},
		},
		{
			Name:    "set",
			Aliases: []string{"s"},
			Usage:   "set key/value, e.g. `confd-cli set key1 value1`",
			//Flags:   app.Flags,
			Action: func(c *cli.Context) {
				set(c)
			},
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "get value of the given key, e.g. `confd-cli get key1`",
			Flags:   app.Flags,
			Action: func(c *cli.Context) {
				get(c)
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "delete the given key, e.g. `confd-cli delete key1`",
			Flags:   app.Flags,
			Action: func(c *cli.Context) {
				delete(c)
			},
		},
	}

	app.Run(os.Args)
}

func printGlobalOptions(c *cli.Context){
	fmt.Println("\nThe input flag args is:")

	for i, k := range c.GlobalFlagNames() {
		fmt.Printf("%d %s:%v\n", i, k, c.GlobalString(k))
	}
}

func printCliConfig(cliConfig *mycli.CliConfig){
	fmt.Println("\nThe cliConfig is:")
	fmt.Printf("%+v\n\n", *cliConfig)
}

func newMyCli(c *cli.Context)(error, mycli.Cli){
	log.SetLevel(log.WarnLevel)
	printGlobalOptions(c)

	var cliConfig *mycli.CliConfig = &mycli.CliConfig{}
	config.InitCliConfig(c.GlobalString("store-type"), c.GlobalString("confdir"), c.GlobalString("config-file"), c.GlobalString("connect-addr"), cliConfig)
	printCliConfig(cliConfig)

	return mycli.New(cliConfig)
}

func getall(c *cli.Context){
	err, mycli := newMyCli(c)
	if err != nil {
		log.Error("New confd cli error: ", err.Error())
	}

	data, err := mycli.GetAll()
	if err!=nil{
		fmt.Println("error when get all key/value")
		return
	}

	fmt.Println("The stored k/v:")

	for k,v:=range data{
		fmt.Printf("%s\t %v \n", k,v)
	}
}

func get(c *cli.Context){
	err, mycli := newMyCli(c)
	if err != nil {
		log.Error("New confd cli error: ", err.Error())
	}

	key := c.Args().Get(0)
	data, err := mycli.GetValue(key)
	if err!=nil{
		fmt.Printf("\nerror when get key[%s], error: %v\n", key,err)
		return
	}

	fmt.Printf("The value of key[%s] is: %v\n",key, data)


}

//confd-cli set key value
func set(c *cli.Context){
	err, mycli := newMyCli(c)
	if err != nil {
		log.Error("New confd cli error: ", err.Error())
	}

	key,value:=c.Args().Get(0), c.Args().Get(1)
	err = mycli.SetValue(key,value)
	if err!=nil{
		fmt.Printf("set error: %s", err)
	}else{
		fmt.Println("set ok.")
	}

}

func delete(c *cli.Context){
	err, mycli := newMyCli(c)
	if err != nil {
		log.Error("New confd cli error: ", err.Error())
	}

	key:=c.Args().Get(0)
	err = mycli.DeleteKey(key)
	if err!=nil{
		fmt.Printf("delete key error: %s", err)
	}else{
		fmt.Println("delete ok.")
	}

}

func main() {
	fmt.Println("Confd cli tool...")
	parseArgs()
}
