package gclient

import (
	ct "context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	sc "strconv"
	"time"

	oa "golang.org/x/oauth2"
)

type dataTokenT struct {
	fileName string
	token    *oa.Token
}

type tokenT struct {
	*dataTokenT
}

func NewToken(fileName string) (t tokenT) {
	t.dataTokenT = &dataTokenT{
		fileName: fileName,
	}
	return
}

// Retrieves a token from a local file.
func (t tokenT) Read() (e error) {
	fmt.Println("Reading token from:", t.fileName)
	file, e := Open(t.fileName, ReadFlag)
	if e != nil {
		return
	}
	defer file.Close()

	t.token = &oa.Token{}
	e = json.NewDecoder(file).Decode(t.token)
	if e != nil {
		log.Printf("Error decoding token file: %v", e)
		return
	}

	return
}

// Saves a token to a local file.
func (t tokenT) Save() (e error) {
	fmt.Println("Saving token to:", t.fileName)
	file, e := Open(t.fileName, WriteFlag)
	if e != nil {
		return
	}
	defer file.Close()

	e = json.NewEncoder(file).Encode(t.token)
	if e != nil {
		log.Printf("Unable to save token: %v", e)
	}

	return
}

func (t tokenT) Archive() (e error) {
	archivedName := t.fileName + ".arch." + sc.FormatInt(time.Now().UnixNano(), 10)
	fmt.Println("Archiving", t.fileName, "as", archivedName)
	e = os.Rename(t.fileName, archivedName)
	if e != nil {
		log.Printf("Unable to archive %s as %s: %v", t.fileName, archivedName, e)
		return
	}

	return
}

// Request a token from the web.
func (t tokenT) Request(ctx ct.Context, conf configT, authCode string) (e error) {
	ctx, cancel := ct.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	t.token, e = conf.Exchange(ctx, authCode)
	if e != nil {
		log.Fatalf("Unable to retrieve token: %v", e)
		return
	}

	return
}

func (t tokenT) NewHttpClient(ctx ct.Context, conf configT) (_ *http.Client) {
	e := t.Read()
	if e != nil {
		fmt.Printf("No token found: %v\n", e)
		authCode, _ := conf.getAuthCode()
		_ = t.Request(ctx, conf, authCode)
		_ = t.Save()
	}

	return conf.Client(ctx, t.token)
}
