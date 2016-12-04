package main

import (
	"bolt"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/urfave/cli"
	. "github.com/yehohanan7/remem/commands"
)

func readText(reader *bufio.Reader) string {
	var text string
	input, err := reader.ReadString('\n')
	for ; input != ".\n" && err == nil; input, err = reader.ReadString('\n') {
		text = text + input
	}
	return text
}

func flags() []cli.Flag {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return []cli.Flag{
		cli.StringFlag{
			Name:  "db",
			Value: usr.HomeDir + "/remem.db",
			Usage: "location of the db",
		},
	}
}

func main() {
	flags := flags()
	app := cli.NewApp()
	app.Version = "1.0"
	app.Name = "remem"

	app.Commands = []cli.Command{
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "show a random note",
			Flags:   flags,
			Action: func(c *cli.Context) {
				Execute(c, func(bucket *bolt.Bucket) {
					fmt.Println(GetFortune(bucket))
				})
			},
		},
		{
			Name:    "add",
			Usage:   "add a note",
			Aliases: []string{"a"},
			Flags:   flags,
			Action: func(c *cli.Context) {
				Execute(c, func(bucket *bolt.Bucket) {
					AddFortune(bucket, readText(bufio.NewReader(os.Stdin)))
				})
			},
		},
	}

	app.Run(os.Args)
}
