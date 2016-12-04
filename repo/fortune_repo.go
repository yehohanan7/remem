package repo

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

const BUCKET_NAME = "fortunes"

func random(n int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(n)
}

type FortuneRepo interface {
	Add(fortune string)
	GetRandom() string
	Close()
}

type BoltFortuneRepo struct {
	db *bolt.DB
}

func (repo *BoltFortuneRepo) Add(fortune string) {
	repo.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))
		key := "0"
		v := string(bucket.Get([]byte("keys")))
		if v != "" {
			lastKey, _ := strconv.Atoi(v)
			key = strconv.Itoa(lastKey + 1)
		}
		_ = bucket.Put([]byte("keys"), []byte(key))
		_ = bucket.Put([]byte(key), []byte(fortune))

		return nil
	})
}

func (repo *BoltFortuneRepo) GetRandom() (fortune string) {
	repo.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))
		key, _ := strconv.Atoi(string(bucket.Get([]byte("keys"))))
		fortune = string(bucket.Get([]byte(strconv.Itoa(random(key + 1)))))
		return nil
	})
	return
}

func (repo *BoltFortuneRepo) Close() {
	repo.db.Close()
}

func newBoltFortuneRepo(location string) (repo *BoltFortuneRepo) {
	db, err := bolt.Open(location, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))
		if err != nil {
			return fmt.Errorf("error initializing the db", err)
		}
		return nil
	})
	return &BoltFortuneRepo{db}
}

func NewFortuneRepo(location string) FortuneRepo {
	return newBoltFortuneRepo(location)
}
