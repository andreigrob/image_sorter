package gclient

import (
	"fmt"
	"net/http"

	ct "context"
)

type DataGClientT struct {
	ctx         ct.Context
	credentials credentialsT
	token       tokenT
	scope       scopeT

	HttpCl *http.Client
}

type GClientT struct {
	*DataGClientT
}

// Creates a Google Drive client.
func NewGClient(ctx ct.Context, credentials string, tokenName string, scope string) (g GClientT) {
	g.DataGClientT = &DataGClientT{
		ctx:         ctx,
		credentials: credentialsT(credentials),
		token:       NewToken(tokenName),
		scope:       scopeT(scope),
	}

	return
}

func (g GClientT) NewHttpClient() (_ *http.Client) {
	conf, _ := g.newConfigT()
	return g.token.NewHttpClient(g.ctx, conf)
}

func (g GClientT) CheckScope() (e error) {
	savedScope, _ := readScope()
	if g.scope == savedScope {
		fmt.Println("Same scope. Can use same", g.token.fileName)
	} else {
		fmt.Println("Different scope.  Archiving", g.token.fileName)
		e = g.token.Archive()
		if e != nil {
			return
		}

		fmt.Println("Saving new scope")
		e = g.scope.save()
		if e != nil {
			fmt.Printf("Unable to save scope: %v\n", e)
		}
	}

	return
}
