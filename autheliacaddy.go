package autheliacaddy

import (
	"context"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

const (

	// cookieName is the name of the cookie that authelia uses.
	cookieName = "authelia"
)

func init() {
	caddy.RegisterModule(Authelia{})
}

// TODO
type Authelia struct{}

// TODO
func (a Authelia) ServeHTTP(writer http.ResponseWriter, request *http.Request, handler caddyhttp.Handler) error {

	// TODO Get hostname from Caddyfile somehow...
	autheliaHostname := ""

	// TODO Create a context...
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Authenticate and authorize the request with Authelia.
	verified, err := verify(ctx, autheliaHostname, http.DefaultClient, request)
	if err != nil {
		return err
	}

	// TODO
	if verified {

	} else {

	}

	return nil
}

// TODO
func (a Authelia) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.micahparks_authelia",
		New: nil, // TODO
	}
}
