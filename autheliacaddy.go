package autheliacaddy

import (
	"context"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// TODO
func init() {
	caddy.RegisterModule(Authelia{})
}

// TODO
type Authelia struct {
	Hostname        string `json:"hostname,omitempty"`
	Timeout         int    `json:"timeout,omitempty"`
	timeoutDuration time.Duration
}

// TODO
func (a Authelia) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.authelia",
		New: nil, // TODO
	}
}

// TODO
func (a *Authelia) Provision(_ caddy.Context) error {

	// If no timeout was given or it was invalid, default to using a one minute timeout.
	if a.Timeout <= 0 {
		a.timeoutDuration = time.Minute
	} else {
		a.timeoutDuration = time.Duration(a.Timeout) * time.Second
	}

	return nil
}

// TODO
func (a Authelia) ServeHTTP(writer http.ResponseWriter, request *http.Request, handler caddyhttp.Handler) error {

	// Create a context for the request to Authelia.
	ctx, cancel := context.WithTimeout(context.Background(), a.timeoutDuration)
	defer cancel()

	// Authenticate and authorize the request with Authelia.
	verified, err := verify(ctx, a.Hostname, http.DefaultClient, request)
	if err != nil {
		return err
	}

	// The request is authenticate and authorized, according to Authelia. Let it through.
	if verified {
		return handler.ServeHTTP(writer, request)
	}

	// Perform a redirect to the Authelia server for authenticate and authorization.
	http.Redirect(writer, request, a.Hostname, http.StatusFound)

	return nil
}
