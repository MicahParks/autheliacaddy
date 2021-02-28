package autheliacaddy

import (
	"context"
	badLogger "log"
	"net/http"
	"net/url"

	"github.com/sanity-io/litter"
)

const (

	// headerHost is the name of the header that contains the host and port of the server in which the original request
	// was sent to.
	headerHost = "Host"

	// headerOriginalURL is the name of the header that contains the originally requested URL.
	headerOriginalURL = "X-Original-URL"
)

var (

	// headersAuthorization holds the set of headers from Authelia that need to be forwarded to the backend for further
	// authorization.
	headersAuthorization = []string{"Remote-Email", "Remote-Groups", "Remote-Name", "Remote-User"}
)

// forwardedHeaders takes the headers from the Authelia response and decides which ones to forward to the backend.
func forwardedHeaders(resp *http.Response) (headers http.Header) {

	// Create the headers to return.
	headers = http.Header{} // TODO Verify initialization.

	// Iterate through the authorization headers to backend.
	for _, key := range headersAuthorization {
		headers.Set(key, resp.Header.Get(key))
	}

	return headers
}

// verify verifies a request with Authelia. If verified, headers will contain the headers to forward with the request to
// the backend.
func (a Authelia) verify(originalReq *http.Request) (verified bool, headers http.Header, err error) {

	// Debug log...
	a.logger.Infow("Performing request to Authelia") // TODO Remove.

	// Clone the original request.
	req := originalReq.Clone(context.Background()) // TODO Verify this.

	// Change the URL of the request so it goes to the Authelia server.
	req.RequestURI = "" // TODO Need to change this to something in verifyURL?
	req.URL = a.verifyURL

	// Parse the original URL's path in relation to the service URL.
	var redirect *url.URL
	if redirect, err = a.serviceURL.Parse(originalReq.URL.Path); err != nil {
		return false, nil, err
	}

	// Set the extra headers for the request.
	req.Header.Set(headerHost, a.verifyURL.Host)
	req.Header.Set(headerOriginalURL, redirect.String())
	req.Host = a.verifyURL.Host

	// Set the redirect for the URL query.
	query := req.URL.Query()
	query.Set("rd", redirect.String())
	req.URL.RawQuery = query.Encode()

	badLogger.Println("doing: " + litter.Sdump(req)) // TODO Remove.

	// Perform the request.
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		panic(err.Error()) // TODO Remove.
		return false, nil, err
	}
	defer resp.Body.Close() // Ignore any error.

	// Decide what to do based on the status code.
	switch resp.StatusCode {
	case 200:
		verified = true
		headers = forwardedHeaders(resp)
	case 401:
		verified = false
	default:
		a.logger.Warnw("Unhandled HTTP status code.",
			"code", resp.StatusCode,
		)
	}

	return verified, headers, nil
}
