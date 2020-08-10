package boltenc

import (
	"fmt"

	"github.com/asdine/storm/v3"
)

// Cryptor is an interface that can be implemented in order
//	to read/write encrypted records to a boltdb database, via storm
type Cryptor interface {
	Decrypt([]byte) ([]byte, error)
	Encrypt([]byte) ([]byte, error)
}

// Record is a struct that is used to encapsulate
//	encrypted credentials in and out of a storm database (boltdb)
type Record struct {
	ID      int `storm:"id,increment"`
	Payload []byte
}

// NewRecord returns an empty Record struct
func NewRecord() *Record {
	return &Record{}
}

// Get retrieves an entry from a storm db, decrypts it, and then returns it
func Get(c Cryptor, db *storm.DB, id int) (map[int][]byte, error) {
	data := map[int][]byte{}
	r := NewRecord()
	err := db.One("ID", id, r)
	if err != nil {
		return data, err
	}

	d, err := c.Decrypt(r.Payload)
	if err != nil {
		return data, err
	}
	data[r.ID] = d

	return data, nil
}

// GetAll retrieves all records
func GetAll(c Cryptor, db *storm.DB) (map[int][]byte, error) {
	var records []Record
	err := db.All(&records)
	if err != nil {
		return map[int][]byte{}, err
	}

	data := map[int][]byte{}
	for _, r := range records {
		d, err := c.Decrypt(r.Payload)
		if err != nil {
			return data, fmt.Errorf("decrypting record, %s", err)
		}
		data[r.ID] = d
	}
	return data, nil
}

// Put appends a record into the database
func Put(c Cryptor, db *storm.DB, payload []byte) (int, error) {
	r := NewRecord()
	cred, err := c.Encrypt(payload)
	if err != nil {
		return -1, err
	}
	r.Payload = cred

	err = db.Save(r)
	if err != nil {
		return -1, err
	}
	return r.ID, nil
}

// Set modifies an existing record
func Set(c Cryptor, db *storm.DB, id int, payload []byte) error {
	newPayload, err := c.Encrypt(payload)
	if err != nil {
		return err
	}

	err = db.UpdateField(&Record{
		ID: id,
	}, "Payload", newPayload)
	if err != nil {
		return err
	}
	return nil
}
