package gclient

import (
	"fmt"
	"log"
)

type scopeT string

const scopeFileName = `scope.txt`

func readScope() (scope scopeT, e error) {
	file, e := Open(scopeFileName, ReadFlag)
	if e != nil {
		return
	}
	defer file.Close()

	// read a line
	_, e = fmt.Fscanln(file, &scope)
	if e != nil {
		log.Printf("Unable to read %s: %v\n", scopeFileName, e)
		return
	}

	return
}

func (s scopeT) save() (e error) {
	file, e := Open(scopeFileName, WriteFlag)
	if e != nil {
		return
	}
	defer file.Close()

	_, e = file.Write([]byte(s + "\n"))
	if e != nil {
		log.Printf("Unable to write %s: %v\n", scopeFileName, e)
		return
	}

	return
}
