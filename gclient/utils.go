package gclient

import (
	"log"
	"os"
)

const ReadFlag = os.O_RDONLY
const WriteFlag = os.O_RDWR | os.O_CREATE | os.O_TRUNC

func Open(fileName string, flags int) (file *os.File, e error) {
	file, e = os.OpenFile(fileName, flags, 0600)
	if e != nil {
		log.Printf("Unable to open %s: %v", fileName, e)
		return
	}

	return
}

func First[T any](a T, _ ...any) (_ T) {
	return a
}

func Second[T any](_ any, b T, _ ...any) (_ T) {
	return b
}
