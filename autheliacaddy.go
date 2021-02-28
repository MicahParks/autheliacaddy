package autheliacaddy

import (
	"context"
	"errors"
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

	// Verify the cookie with Authelia.
	var verified bool
	if verified, err = verify(ctx, autheliaHostname, http.DefaultClient, request); err != nil {
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
