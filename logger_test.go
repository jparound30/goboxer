package goboxer

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// captureLog returns a JSON slog logger writing to buf, and a function that parses the records.
func captureLog(t *testing.T) (*slog.Logger, func() []map[string]any) {
	t.Helper()
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger, func() []map[string]any {
		var records []map[string]any
		for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
			if line == "" {
				continue
			}
			var rec map[string]any
			if err := json.Unmarshal([]byte(line), &rec); err != nil {
				t.Fatalf("log line is not JSON: %s", line)
			}
			records = append(records, rec)
		}
		return records
	}
}

func findRecord(records []map[string]any, msg string) map[string]any {
	for _, r := range records {
		if r["msg"] == msg {
			return r
		}
	}
	return nil
}

func TestLog_RequestMasksSecrets(t *testing.T) {
	logger, records := captureLog(t)
	Log, LogOpts = logger, LogOptions{RequestBody: true}
	defer func() { Log, LogOpts = nil, LogOptions{} }()

	form := "grant_type=refresh_token&refresh_token=RT_SECRET&client_id=ID&client_secret=CS_SECRET"
	req, _ := http.NewRequest(http.MethodPost, "https://api.box.com/oauth2/token", strings.NewReader(form))
	req.Header.Set(httpHeaderAuthorization, "Bearer ACCESS_SECRET")
	req.Header.Set(httpHeaderContentType, ContentTypeFormUrlEncoded)

	logRequest(http.MethodPost, req)

	rec := findRecord(records(), "goboxer: request")
	if rec == nil {
		t.Fatalf("request record not logged")
	}
	raw, _ := json.Marshal(rec)
	for _, secret := range []string{"ACCESS_SECRET", "RT_SECRET", "CS_SECRET"} {
		if strings.Contains(string(raw), secret) {
			t.Errorf("log record leaks %q: %s", secret, raw)
		}
	}
	headers := rec["headers"].(map[string]any)
	if got := headers["Authorization"].([]any)[0]; got != "Bearer ***" {
		t.Errorf("Authorization logged as %v, want \"Bearer ***\"", got)
	}
	body := rec["body"].(string)
	if !strings.Contains(body, "client_id=ID") || !strings.Contains(body, "grant_type=refresh_token") {
		t.Errorf("non-secret form fields should be kept: %s", body)
	}
}

func TestLog_BodiesOnlyWhenEnabled(t *testing.T) {
	logger, records := captureLog(t)
	Log = logger
	defer func() { Log, LogOpts = nil, LogOptions{} }()

	req, _ := http.NewRequest(http.MethodPost, "https://api.box.com/2.0/folders", strings.NewReader(`{"name":"x"}`))
	req.Header.Set(httpHeaderContentType, ContentTypeApplicationJson)

	LogOpts = LogOptions{}
	logRequest(http.MethodPost, req)
	logResponseBody(ContentTypeApplicationJson, []byte(`{"id":"1"}`))

	LogOpts = LogOptions{RequestBody: true, ResponseBody: true}
	logRequest(http.MethodPost, req)
	logResponseBody(ContentTypeApplicationJson, []byte(`{"id":"1"}`))
	logResponseBody("application/octet-stream", []byte("binary"))

	recs := records()
	var requestBodies, responseBodies int
	for _, r := range recs {
		switch r["msg"] {
		case "goboxer: request":
			if _, ok := r["body"]; ok {
				requestBodies++
			}
		case "goboxer: response body":
			responseBodies++
		}
	}
	if requestBodies != 1 {
		t.Errorf("request bodies logged = %d, want 1 (only when enabled)", requestBodies)
	}
	if responseBodies != 1 {
		t.Errorf("response bodies logged = %d, want 1 (only when enabled and JSON)", responseBodies)
	}
}

func TestLog_RetryRecords(t *testing.T) {
	logger, records := captureLog(t)
	Log = logger
	defer func() { Log = nil }()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer ts.Close()

	apiConn := commonInit(ts.URL)
	apiConn.MaxRequestAttempts = 2
	if _, err := NewRequest(apiConn, apiConn.BaseURL+"users/me", GET, nil, nil).Send(); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	recs := records()
	retrying := findRecord(recs, "goboxer: retrying")
	if retrying == nil || retrying["level"] != "INFO" || retrying["attempt"] != float64(1) || retrying["status"] != float64(429) {
		t.Errorf("retrying record = %v", retrying)
	}
	limit := findRecord(recs, "goboxer: retry limit reached")
	if limit == nil || limit["level"] != "WARN" || limit["attempts"] != float64(2) {
		t.Errorf("retry limit record = %v", limit)
	}
	if findRecord(recs, "goboxer: response") == nil || findRecord(recs, "goboxer: request") == nil {
		t.Errorf("request/response summaries should be logged at debug level")
	}
}

func TestLog_NilDisablesLogging(t *testing.T) {
	Log = nil
	req, _ := http.NewRequest(http.MethodGet, "https://api.box.com/2.0/users/me", nil)
	// must not panic
	logRequest(http.MethodGet, req)
	logResponse(&http.Response{StatusCode: 200, Header: http.Header{}}, 1)
	logResponseBody(ContentTypeApplicationJson, []byte(`{}`))
}
