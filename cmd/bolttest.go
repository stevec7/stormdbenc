package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/asdine/storm/v3"
	"github.com/stevec7/sshencode/pkg/sshencode"

	//"github.com/stevec7/random/boltdbenc/pkg/boltenc"
	"github.com/stevec7/random/boltdbenc/pkg/boltenc"
)

/*
type Cred struct {
	ID      int `storm:"id,increment"`
	Payload []byte
}

type secret string

func (s secret) String() string {
	return fmt.Sprintf("REDACTED")
}

type cred struct {
	Username string
	Host     string
	Password secret
}

func NewCred(user, host, password string) *cred {
	return &cred{
		Username: user,
		Host:     host,
		Password: secret(password),
	}
}

*/

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

	//err = addCred(db, uuid.New().String(), []byte(cipherTexts[0]))
	c := cred{
		Username: "admin",
		Password: "P@ssw0rd",
		Host:     "localhost",
	}

	cr := NewC()
	cr.a = agent

	record, err := boltenc.Put(cr, db, c.Bytes())
	if err != nil {
		log.Fatalf("Error getting record, %s", err)
	}

	fmt.Printf("Record: %+v\n", record)

	/*
		err = boltenc.Set(cr, db, 1, c.Bytes())
		if err != nil {
			log.Fatalf("Error getting record, %s", err)
		}
	*/

	//fmt.Printf("Record: %+v\n", record)

	results, err := boltenc.GetAll(cr, db)
	if err != nil {
		log.Fatalf("Error getting all entries, %s", err)
	}

	for i, r := range results {
		fmt.Printf("i: %d, r: %+v\n", i, string(r))
	}

	/*
		enc, err := agent.Encrypt(c.Bytes())
		if err != nil {
			log.Fatalf("Error %s", err)
		}
		credential := Cred{Payload: enc}
		err = db.Save(&credential)
		if err != nil {
			log.Fatalf("error saving, %s", err)
		}
		id := credential.ID
		fmt.Printf("credential struct: %+v\n", credential)
		var entry Cred
		err = db.One("ID", id, &entry)

		decrypted, err := agent.Decrypt(entry.Payload)
		if err != nil {
			log.Fatalf("Error %s", err)
		}
		fmt.Printf("d: %+v\n", string(decrypted))
		var cc cred
		err = json.Unmarshal(decrypted, &cc)
		if err != nil {
			log.Fatalf("Error %s", err)
		}
		fmt.Printf("struct: %+v\n", cc)
	*/

}
