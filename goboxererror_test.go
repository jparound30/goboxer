package goboxer

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestApiOtherError_Formatting(t *testing.T) {
	cause := errors.New("connection reset")
	wrapped := fmt.Errorf("failed to send request: %w", cause)

	tests := []struct {
		name string
		err  error
		want string
	}{
		{"without body", newApiOtherError(wrapped, ""), "failed to send request: connection reset\n"},
		{"with body", newApiOtherError(wrapped, `{"x":1}`), "failed to send request: connection reset\nbox api response:\n-----\n{\"x\":1}\n-----\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
			// every fmt verb prints the same message; there is no stack frame output any more
			for _, verb := range []string{"%v", "%+v", "%s"} {
				if got := fmt.Sprintf(verb, tt.err); got != tt.want {
					t.Errorf("Sprintf(%s) = %q, want %q", verb, got, tt.want)
				}
			}
			if got := fmt.Sprintf("%q", tt.err); got != fmt.Sprintf("%q", tt.want) {
				t.Errorf("Sprintf(%%q) = %s, want %q", got, tt.want)
			}

			if !errors.Is(tt.err, cause) {
				t.Errorf("errors.Is() should find the root cause through ApiOtherError")
			}
			var target *ApiOtherError
			if !errors.As(tt.err, &target) {
				t.Fatalf("errors.As() should match *ApiOtherError")
			}
			if target.Unwrap() != wrapped {
				t.Errorf("Unwrap() = %v, want %v", target.Unwrap(), wrapped)
			}
		})
	}
}

func TestNewApiStatusError(t *testing.T) {
	body := []byte(`{"type":"error","status":404,"code":"not_found","message":"Not Found","request_id":"abc123","context_info":{"errors":[]}}`)

	err := NewApiStatusError(body)
	var statusErr *ApiStatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("NewApiStatusError() = %T, want *ApiStatusError", err)
	}
	if statusErr.Status != 404 || statusErr.Code != "not_found" || statusErr.RequestId != "abc123" {
		t.Errorf("parsed fields = %+v", statusErr)
	}
	want := "HTTP status code: [404], Message: [Not Found], RequestId: [abc123], ContextInfo: [map[errors:[]]]"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	for _, verb := range []string{"%v", "%+v", "%s"} {
		if got := fmt.Sprintf(verb, err); got != want {
			t.Errorf("Sprintf(%s) = %q, want %q", verb, got, want)
		}
	}
	if statusErr.Unwrap() != nil {
		t.Errorf("Unwrap() = %v, want nil", statusErr.Unwrap())
	}
	if errors.Unwrap(err) != nil {
		t.Errorf("errors.Unwrap() = %v, want nil", errors.Unwrap(err))
	}
}

func TestNewApiStatusError_InvalidBody(t *testing.T) {
	err := NewApiStatusError([]byte(`<html>502 Bad Gateway</html>`))
	var otherErr *ApiOtherError
	if !errors.As(err, &otherErr) {
		t.Fatalf("NewApiStatusError() with invalid body = %T, want *ApiOtherError", err)
	}
	msg := err.Error()
	if !strings.Contains(msg, "response json marshaling error") {
		t.Errorf("Error() = %q, want json error description", msg)
	}
	if !strings.Contains(msg, "<html>502 Bad Gateway</html>") {
		t.Errorf("Error() = %q, want the raw body included", msg)
	}
}
