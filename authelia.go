package autheliacaddy

import (
	"context"
	"net/http"
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
func (a Authelia) verify(ctx context.Context, originalReq *http.Request) (verified bool, headers http.Header, err error) {

	// Clone the original request.
	req := originalReq.Clone(ctx) // TODO Verify this

	// Set the extra headers for the request.
	//
	// TODO Verify.
	req.Header.Set(headerHost, a.url.Host)
	req.Header.Set(headerOriginalURL, req.URL.String())

	// Change the URL of the request so it goes to the Authelia server.
	req.URL = a.url

	// Perform the request.
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
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
