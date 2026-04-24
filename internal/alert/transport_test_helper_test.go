package alert

import (
	"net/http"
	"net/url"
	"strings"
)

// rewriteTransport returns an http.RoundTripper that redirects all requests
// to the given base URL, allowing tests to intercept HTTP calls.
type urlRewriteTransport struct {
	baseURL string
}

func rewriteTransport(baseURL string) http.RoundTripper {
	return &urlRewriteTransport{baseURL: baseURL}
}

func (t *urlRewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base, err := url.Parse(t.baseURL)
	if err != nil {
		return nil, err
	}
	req = req.Clone(req.Context())
	req.URL.Scheme = base.Scheme
	req.URL.Host = base.Host
	// Strip any path prefix mismatch by keeping only the last path segment
	// so test servers receive requests at their root.
	parts := strings.Split(req.URL.Path, "/")
	if len(parts) > 0 {
		req.URL.Path = "/" + parts[len(parts)-1]
	}
	return http.DefaultTransport.RoundTrip(req)
}
