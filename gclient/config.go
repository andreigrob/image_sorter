package gclient

import (
	"fmt"
	"log"

	oa "golang.org/x/oauth2"
	gg "golang.org/x/oauth2/google"
)


type configT struct {
	*oa.Config
}

func (g GClientT) newConfigT() (c configT, e error) {
	creds, _ := g.credentials.read()

	c.Config, e = gg.ConfigFromJSON(creds, string(g.scope))
	if e != nil {
		log.Fatalf("Unable to parse %s: %v", g.credentials, e)
		return
	}

	return
}

// Retrieves an authorization code from Google OAuth.
func (c configT) getAuthCode() (authCode string, e error) {
	fmt.Printf("Go to this url:\n%v\nAuthorization code: ",
		c.AuthCodeURL(`state-token`, oa.AccessTypeOffline),
	)

	if _, e = fmt.Scanln(&authCode); e != nil {
		log.Fatalf("Unable to read authorization code: %v", e)
		return
	}

	return
}
