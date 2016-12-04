package commands

import (
	"bolt"
	"strconv"
)

func GetFortune(bucket *bolt.Bucket) string {
	key, _ := strconv.Atoi(string(bucket.Get([]byte("keys"))))
	return string(bucket.Get([]byte(strconv.Itoa(random(key + 1)))))
}
