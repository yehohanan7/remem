package commands

import (
	"bolt"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

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

func Execute(c *cli.Context, fn func(*bolt.Bucket)) {
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
