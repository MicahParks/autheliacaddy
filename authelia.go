package autheliacaddy

import (
	"context"
	"net/http"
	"net/url"
)

const (

	// endpointVerify is the path to the endpoint that will verify incoming HTTP requests.
	endpointVerify = "/verify"

	// headerHost is the name of the header that contains the host and port of the server in which the original request
	// was sent to.
	headerHost = "Host"

	// headerOriginalURL is the name of the header that contains the originally requested URL.
	headerOriginalURL = "X-Original-URL"
)

// TODO
// hostname should not have an http prefix.
func verify(ctx context.Context, autheliaHostname string, client *http.Client, originalReq *http.Request) (verified bool, err error) {

	// Create the full URL to make a request to.
	//
	// TODO HTTPS?
	u := "http://" + autheliaHostname + endpointVerify

	// Clone the original request.
	req := originalReq.Clone(ctx) // TODO Verify this

	// Set the extra headers for the request.
	//
	// TODO Verify.
	req.Header.Set(headerHost, autheliaHostname)
	req.Header.Set(headerOriginalURL, req.URL.String())

	// Change the URL of the request so it goes to the Authelia server.
	if req.URL, err = url.Parse(u); err != nil {
		return false, err
	}

	// Perform the request.
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return false, err
	}
	defer resp.Body.Close() // Ignore any error.

	// Decide what to do based on the status code.
	switch resp.StatusCode {
	case 200:
		verified = true
	case 401:
		verified = false
	default:
		// TODO
	}

	return verified, nil
}
