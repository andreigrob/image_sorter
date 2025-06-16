package gclient

import (
	"log"
	"os"
)

type credentialsT string

func (c credentialsT) read() (creds []byte, e error) {
	creds, e = os.ReadFile(string(c))
	if e != nil {
		log.Fatalf("Unable to read credentials from %s: %v", c, e)
		return
	}

	return
}