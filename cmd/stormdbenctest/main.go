package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/asdine/storm/v3"
	"github.com/stevec7/sshencode/pkg/sshencode"
	"github.com/stevec7/stormdbenc/pkg/stormdbenc"
)


type secret string

type cred struct {
	Username string
	Host     string
	Password secret
}

func (c *cred) Bytes() []byte {
	r, _ := json.Marshal(c)
	return r
}

type C struct {
	a *sshencode.Agent
}

func (c *C) Decrypt(b []byte) ([]byte, error) {
	d, err := c.a.Decrypt(b)
	if err != nil {
		return []byte{}, err
	}
	return d, nil
}

func (c *C) Encrypt(b []byte) ([]byte, error) {
	d, err := c.a.Encrypt(b)
	if err != nil {
		return []byte{}, err
	}
	return d, nil
}

func NewC() *C {
	return &C{}
}

func main() {
	filename := flag.String("filename", "", "database file")
	flag.Parse()
	prefix := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
	agent, err := sshencode.Configure(prefix)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		os.Exit(1)
	}

	db, err := storm.Open(*filename)
	if err != nil {
		log.Fatalf("error opening, %s", err)
	}
	defer db.Close()

	c := cred{
		Username: "admin",
		Password: "P@ssw0rd",
		Host:     "localhost",
	}

	cr := NewC()
	cr.a = agent

	record, err := stormdbenc.Put(cr, db, c.Bytes())
	if err != nil {
		log.Fatalf("Error getting record, %s", err)
	}

	fmt.Printf("Record: %+v\n", record)

	results, err := stormdbenc.GetAll(cr, db)
	if err != nil {
		log.Fatalf("Error getting all entries, %s", err)
	}

	for i, r := range results {
		fmt.Printf("i: %d, r: %+v\n", i, string(r))
	}
}
