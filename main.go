package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"strconv"

	"math/rand"

	"github.com/boltdb/bolt"
	"github.com/urfave/cli"
)

const (
	BUCKET_NAME = "Fortunes"
	DEFAULT_KEY = "0"
)

func random(n int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(n)
}

func getNewKey(bucket *bolt.Bucket) (key string) {
	v := string(bucket.Get([]byte("keys")))
	if v == "" {
		key = DEFAULT_KEY
	} else {
		lastKey, _ := strconv.Atoi(v)
		key = strconv.Itoa(lastKey + 1)
	}
	_ = bucket.Put([]byte("keys"), []byte(key))
	return
}

func addFortune(bucket *bolt.Bucket) {
	text := readText(bufio.NewReader(os.Stdin))
	_ = bucket.Put([]byte(getNewKey(bucket)), []byte(text))
}

func printFortune(bucket *bolt.Bucket) {
	key, _ := strconv.Atoi(string(bucket.Get([]byte("keys"))))
	v := string(bucket.Get([]byte(strconv.Itoa(random(key + 1)))))
	fmt.Println(v)
}

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

func execute(c *cli.Context, fn func(*bolt.Bucket)) {
	db, err := bolt.Open(c.String("db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))
		if err != nil {
			return fmt.Errorf("error initializing the db", err)
		}
		fn(bucket)
		return nil
	})
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
				execute(c, printFortune)
			},
		},
		{
			Name:    "add",
			Usage:   "add a note",
			Aliases: []string{"a"},
			Flags:   flags,
			Action: func(c *cli.Context) {
				execute(c, addFortune)
			},
		},
	}

	app.Run(os.Args)
}
