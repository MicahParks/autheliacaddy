package autheliacaddy

import (
	"bytes"
	"context"
	"errors"
	"net/http"
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

	// Create the request to verify
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, u, bytes.NewReader(nil)); err != nil {
		return false, err
	}

	// Set the headers for the request.
	//
	// TODO Verify.
	req.Header.Set(headerHost, autheliaHostname)
	req.Header.Set(headerOriginalURL, originalURL)

	// Copy any authentication cookie from the header.
	cookie, err := originalReq.Cookie(cookieName)
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			return false, err
		}
	}
	if cookie != nil {
		req.AddCookie(cookie)
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
