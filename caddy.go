package autheliacaddy

import (
	"fmt"
	badLogger "log"
	"net/http"
	"net/url"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/prometheus/common/log"
	"github.com/sanity-io/litter"
	"go.uber.org/zap"
)

// Interface guards.
var (
	_ caddy.Provisioner           = (*Authelia)(nil)
	_ caddyfile.Unmarshaler       = (*Authelia)(nil)
	_ caddyhttp.MiddlewareHandler = (*Authelia)(nil)
)

// init is apart of creating a Caddy v2 module.
func init() {
	caddy.RegisterModule(Authelia{})
	httpcaddyfile.RegisterHandlerDirective("authelia", parseCaddyfileHandler)
}

// Authelia is a Caddy v2 module that will perform authentication and authorization of requests with a Authelia
// instance.
type Authelia struct {
	ServiceURL string `json:"service_url,omitempty"`
	VerifyURL  string `json:"verify_url,omitempty"`
	logger     *zap.SugaredLogger
	serviceURL *url.URL
	verifyURL  *url.URL
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

	// Turn the raw verification URL into the correct Go type.
	var err error
	if a.verifyURL, err = url.Parse(a.VerifyURL); err != nil {
		return fmt.Errorf("failed to parse Authelia verification URL: %w", err)
	}

	// Turn the raw service URL into the correct Go type.
	if a.serviceURL, err = url.Parse(a.ServiceURL); err != nil {
		return fmt.Errorf("failed to parse service (protected by Authelia) verification URL: %w", err)
	}

	return nil
}

// ServeHTTP implements the caddyhttp.MiddlewareHandler interface. It serves as an HTTP middleware to authenticate
// requests to Authelia.
func (a Authelia) ServeHTTP(writer http.ResponseWriter, request *http.Request, handler caddyhttp.Handler) error {

	badLogger.Println(litter.Sdump(request)) // TODO Remove.

	//// Do not match any requests to the Authelia server to prevent a loop.
	//if request.URL.Host != a.url.Host {
	//	a.logger.Infow("Not sending request to Authelia.",
	//		"host", request.URL.Host,
	//		"url", request.URL.String(),
	//	) // TODO Remove.
	//	return handler.ServeHTTP(writer, request)
	//}

	a.logger.Infow("Sending request to Authelia.",
		"host", request.URL.Host,
		"url", request.URL.String(),
	) // TODO Remove.

	// Authenticate and authorize the request with Authelia.
	verified, headers, err := a.verify(request)
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
//     authelia <verify_url> <service_url>
//
func (a *Authelia) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {

	// Iterate through the tokens.
	for d.Next() {

		// Get all of the arguments.
		arguments := d.RemainingArgs()

		log.Infof("arguments: %v", arguments) // TODO Remove.

		// Confirm all three arguments are present.
		if len(arguments) != 2 {
			return d.ArgErr()
		}

		// Assign the arguments to the data structure.
		a.VerifyURL = arguments[0]
		a.ServiceURL = arguments[1]
	}

	return nil
}

// parseCaddyfileHandler unmarshals tokens from h into a new middleware handler value.
func parseCaddyfileHandler(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var a Authelia
	err := a.UnmarshalCaddyfile(h.Dispenser)
	return a, err
}
