package bolt

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type Cred struct {
	Username string
	Host     string
	Password string
}

func setupDB(filename string) (*bolt.DB, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("CREDS"))
		if err != nil {
			return fmt.Errorf("could not create creds bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	fmt.Println("DB Setup Done")
	return db, nil
}
