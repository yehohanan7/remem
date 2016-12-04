package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/urfave/cli"
	"github.com/yehohanan7/remem/repo"
)

func readText(reader *bufio.Reader, terminator string) string {
	var text string
	input, err := reader.ReadString('\n')
	for ; input != terminator && err == nil; input, err = reader.ReadString('\n') {
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

func addFortune(c *cli.Context) {
	repo := repo.NewFortuneRepo(c.String("db"))
	defer repo.Close()
	repo.Add(readText(bufio.NewReader(os.Stdin), ".\n"))
}

func printFortune(c *cli.Context) {
	repo := repo.NewFortuneRepo(c.String("db"))
	defer repo.Close()
	fmt.Println(repo.GetRandom())
}

func main() {
	flags := flags()

	app := cli.NewApp()
	app.Version = "1.0"
	app.Name = "remem"
	app.Flags = flags
	app.Commands = []cli.Command{
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "show a random note",
			Flags:   flags,
			Action:  printFortune,
		},
		{
			Name:    "add",
			Usage:   "add a note",
			Aliases: []string{"a"},
			Flags:   flags,
			Action:  addFortune,
		},
	}
	app.Action = addFortune
	app.Run(os.Args)
}
