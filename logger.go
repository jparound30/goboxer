package goboxer

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

// Log is the structured logger used by the library. It is nil by default, which disables logging.
//
// Request and response summaries are logged at slog.LevelDebug, retries at slog.LevelInfo and
// failed attempts at slog.LevelWarn. Set it to any *slog.Logger, for example:
//
//	goboxer.Log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
//
// The Authorization header and secrets in form bodies (client_secret, refresh_token, ...) are masked.
var Log *slog.Logger

// LogOptions selects the optional details that are logged in addition to request and response summaries.
type LogOptions struct {
	// RequestBody logs JSON and form encoded request bodies. Secrets in form bodies are masked.
	RequestBody bool
	// ResponseBody logs JSON response bodies once they have been read.
	ResponseBody bool
}

// LogOpts configures what Log records. Bodies are not logged by default.
var LogOpts LogOptions

const redactedValue = "***"

var sensitiveFormKeys = map[string]bool{
	"client_secret": true,
	"refresh_token": true,
	"access_token":  true,
	"assertion":     true,
	"code":          true,
}

// redactHeaders returns a copy of h with the Authorization header value masked.
func redactHeaders(h http.Header) http.Header {
	out := make(http.Header, len(h))
	for k, v := range h {
		if http.CanonicalHeaderKey(k) == httpHeaderAuthorization {
			masked := make([]string, len(v))
			for i, s := range v {
				if scheme, _, ok := strings.Cut(s, " "); ok {
					masked[i] = scheme + " " + redactedValue
				} else {
					masked[i] = redactedValue
				}
			}
			out[k] = masked
			continue
		}
		out[k] = append([]string(nil), v...)
	}
	return out
}

// redactFormBody masks sensitive values in an application/x-www-form-urlencoded body.
// If the body cannot be parsed it is replaced entirely.
func redactFormBody(body string) string {
	values, err := url.ParseQuery(body)
	if err != nil {
		return redactedValue
	}
	for key := range values {
		if sensitiveFormKeys[key] {
			values[key] = []string{redactedValue}
		}
	}
	return values.Encode()
}
