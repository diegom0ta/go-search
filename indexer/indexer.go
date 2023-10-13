package indexer

import (
	"math/rand"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type Indexer struct {
	db *bolt.DB
}

func NewIndexer(db *bolt.DB) *Indexer {
	return &Indexer{
		db: db,
	}
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (i *Indexer) Index(data string) {
	sd := strings.Split(data, " ")
	kv := make(map[string]string)

	for _, d := range sd {
		s := generateRandomString(4)
		kv[s] = d
	}

	// Open a writable transaction
	i.db.Update(func(tx *bolt.Tx) error {
		// For each word, add an entry in the database
		for word, urls := range kv {
			// Use the word as the key and the URLs as the value
			err := tx.Bucket([]byte("IndexBucket")).Put([]byte(word), []byte(urls))
			if err != nil {
				return err
			}
		}
		return nil
	})
}
