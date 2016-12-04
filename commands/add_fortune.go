package commands

import "bolt"

func AddFortune(bucket *bolt.Bucket, fortune string) {
	_ = bucket.Put([]byte(getNewKey(bucket)), []byte(fortune))
}
