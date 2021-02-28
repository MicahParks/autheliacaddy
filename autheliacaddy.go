package autheliacaddy

import (
	"context"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
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

	// The request is authenticate and authorized, according to Authelia. Let it through.
	if verified {
		return handler.ServeHTTP(writer, request)
	}

	// Perform a redirect to the Authelia server for authenticate and authorization.
	http.Redirect(writer, request, autheliaHostname, http.StatusFound)

	return nil
}

// TODO
func (a Authelia) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.micahparks_authelia",
		New: nil, // TODO
	}
}
