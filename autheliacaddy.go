package autheliacaddy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

// Interface guards.
var (
	_ caddy.Provisioner           = (*Authelia)(nil)
	_ caddyfile.Unmarshaler       = (*Authelia)(nil)
	_ caddyhttp.MiddlewareHandler = (*Authelia)(nil)
)

// TODO
func init() {
	caddy.RegisterModule(Authelia{})
	httpcaddyfile.RegisterHandlerDirective("authelia", parseCaddyfileHandler)
}

// TODO
type Authelia struct {
	Prefix     string `json:"prefix"`
	VerifyURL  string `json:"url,omitempty"`
	RawTimeout string `json:"raw_timeout"`
	logger     *zap.SugaredLogger
	timeout    time.Duration
	url        *url.URL
}

// CaddyModule implements the caddy.Module interface. It creates a new Authelia module.
func (a Authelia) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.authelia",
		New: func() caddy.Module { return new(Authelia) },
	}
}

// Provision implements the caddy.Provisioner interface. It does work before the module can be used.
func (a *Authelia) Provision(ctx caddy.Context) error {

	// Add the logger.
	a.logger = ctx.Logger(a).Sugar()

	// Turn the raw URL into the correct Go type.
	var err error
	if a.url, err = url.Parse(a.VerifyURL); err != nil {
		return fmt.Errorf("failed to parse Authelia URL: %w", err)
	}

	// Parse the timeout as an unsigned integer.
	timeout, err := strconv.ParseInt(a.RawTimeout, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse timeout as quatity of seconds: %s", a.RawTimeout)
	}

	// If no timeout was given or it was invalid, default to using a one minute timeout.
	if timeout <= 0 {
		ctx.Logger(a).Sugar().Infow("Given timeout for Authelia was invalid. Defaulting to one minute.",
			"timeout", timeout,
		) // TODO Remove?
		a.timeout = time.Minute
	} else {
		a.timeout = time.Duration(timeout) * time.Second
	}

	return nil
}

// ServeHTTP implements the caddyhttp.MiddlewareHandler interface. It serves as an HTTP middleware to authenticate
// requests to Authelia.
func (a Authelia) ServeHTTP(writer http.ResponseWriter, request *http.Request, handler caddyhttp.Handler) error {

	// Determine if the request has the required prefix.
	if !strings.HasPrefix(request.URL.Path, a.Prefix) {
		return handler.ServeHTTP(writer, request)
	}

	// Create a context for the request to Authelia.
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout*time.Second)
	defer cancel()

	// Authenticate and authorize the request with Authelia.
	verified, headers, err := a.verify(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to verify request with Authelia: %w", err)
	}

	// The request is authenticate and authorized, according to Authelia.
	if verified {

		// Add the forwarded headers to the request.
		for key := range headers {
			request.Header.Set(key, headers.Get(key))
		}

		// Let the request go on to the next handler.
		return handler.ServeHTTP(writer, request)
	}

	// Perform a redirect to the Authelia server for authenticate and authorization.
	http.Redirect(writer, request, a.VerifyURL, http.StatusFound)

	return nil
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler. Syntax:
//
//     authelia <prefix> <verify_url> <timeout>
//
func (a *Authelia) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {

	panic("pull something new")

	// Iterate through the tokens.
	for d.Next() {

		// Get all of the arguments.
		arguments := d.RemainingArgs()

		a.logger.Infow("",
			"arguments", arguments,
		)

		// Confirm all three arguments are present.
		if len(arguments) != 3 {
			return d.ArgErr()
		}

		// Assign the arguments to the data structure.
		a.Prefix = arguments[0]
		a.VerifyURL = arguments[1]
		a.RawTimeout = arguments[2]
	}

	return nil
}

// parseCaddyfileHandler unmarshals tokens from h into a new middleware handler value.
func parseCaddyfileHandler(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var a Authelia
	err := a.UnmarshalCaddyfile(h.Dispenser)
	return a, err
}
