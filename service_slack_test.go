package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type fakeSlackViewOpener struct {
	triggerID string
	view      slackModalView
	err       error
	calls     int
}

func (f *fakeSlackViewOpener) OpenView(ctx context.Context, triggerID string, view slackModalView) error {
	f.triggerID = triggerID
	f.view = view
	f.calls++
	return f.err
}

func TestVerifySlackRequestSignature(t *testing.T) {
	now := time.Unix(1779721200, 0)
	body := []byte("payload=test")
	headers := signedSlackHeaders("slack-secret", body, now)

	cases := []struct {
		name    string
		secret  string
		headers http.Header
		now     time.Time
		wantErr bool
	}{
		{
			name:    "OK",
			secret:  "slack-secret",
			headers: headers,
			now:     now,
		},
		{
			name:    "NG empty secret",
			secret:  "",
			headers: headers,
			now:     now,
			wantErr: true,
		},
		{
			name:    "NG expired timestamp",
			secret:  "slack-secret",
			headers: headers,
			now:     now.Add(slackSignatureMaxAge + time.Second),
			wantErr: true,
		},
		{
			name:   "NG invalid signature",
			secret: "slack-secret",
			headers: http.Header{
				"X-Slack-Request-Timestamp": []string{fmt.Sprintf("%d", now.Unix())},
				"X-Slack-Signature":         []string{"v0=invalid"},
			},
			now:     now,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := verifySlackRequestSignature(c.secret, c.headers, body, c.now)
			if (err != nil) != c.wantErr {
				t.Fatalf("Unexpected error: got=%v wantErr=%t", err, c.wantErr)
			}
		})
	}
}

func TestParseAndVerifySlackActionPayload(t *testing.T) {
	now := time.Unix(1779721200, 0)
	payload := buildTestSlackActionPayload(slackActionPend, 1001, 123456, "action-secret", now)
	rawPayload := mustJSON(t, payload)

	got, err := parseAndVerifySlackActionPayload(rawPayload, "action-secret", now)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got.Action != slackActionPend || got.ProjectID != 1001 || got.FindingID != 123456 {
		t.Fatalf("Unexpected payload: got=%+v", got)
	}

	expired := payload
	expired.ExpiresAt = now.Add(-time.Second).Unix()
	expired.Signature = signSlackActionPayload(expired, "action-secret")
	if _, err := parseAndVerifySlackActionPayload(mustJSON(t, expired), "action-secret", now); err == nil {
		t.Fatal("Expected expired payload to be rejected")
	}

	tampered := payload
	tampered.ProjectID = 2002
	if _, err := parseAndVerifySlackActionPayload(mustJSON(t, tampered), "action-secret", now); err == nil {
		t.Fatal("Expected tampered payload to be rejected")
	}

	invalidAction := payload
	invalidAction.Action = "delete"
	invalidAction.Signature = signSlackActionPayload(invalidAction, "action-secret")
	if _, err := parseAndVerifySlackActionPayload(mustJSON(t, invalidAction), "action-secret", now); err == nil {
		t.Fatal("Expected invalid action to be rejected")
	}
}

func TestSlackActionHandlerOpensPendModal(t *testing.T) {
	now := time.Now()
	opener := &fakeSlackViewOpener{}
	svc := &gatewayService{
		slackSigningSecret:       "slack-secret",
		slackActionSigningSecret: "action-secret",
		slackViewOpener:          opener,
	}
	rawActionPayload := mustJSON(t, buildTestSlackActionPayload(slackActionPend, 1001, 123456, "action-secret", now))
	req := newSignedSlackActionRequest(t, "slack-secret", rawActionPayload, now)
	rec := httptest.NewRecorder()

	svc.slackActionHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Unexpected status: got=%d body=%s", rec.Code, rec.Body.String())
	}
	if opener.calls != 1 {
		t.Fatalf("Unexpected Slack modal call count: got=%d", opener.calls)
	}
	if opener.triggerID != "trigger-1" {
		t.Fatalf("Unexpected trigger ID: got=%s", opener.triggerID)
	}
	if opener.view.CallbackID != slackActionCallbackPend {
		t.Fatalf("Unexpected callback ID: got=%s", opener.view.CallbackID)
	}
	if opener.view.PrivateMetadata != rawActionPayload {
		t.Fatalf("Unexpected private metadata: got=%s", opener.view.PrivateMetadata)
	}
	if opener.view.Title.Text != "Pending Finding" {
		t.Fatalf("Unexpected title: got=%s", opener.view.Title.Text)
	}
	if len(opener.view.Blocks) != 3 {
		t.Fatalf("Unexpected block count: got=%d", len(opener.view.Blocks))
	}
	if opener.view.Blocks[1].BlockID != "pend_deadline" || opener.view.Blocks[1].Element.Type != "datepicker" {
		t.Fatalf("Unexpected PEND deadline block: got=%+v", opener.view.Blocks[1])
	}
	if !opener.view.Blocks[2].Optional {
		t.Fatalf("Expected PEND note to be optional: got=%+v", opener.view.Blocks[2])
	}
}

func TestSlackActionHandlerOpensArchiveModal(t *testing.T) {
	now := time.Now()
	opener := &fakeSlackViewOpener{}
	svc := &gatewayService{
		slackSigningSecret:       "slack-secret",
		slackActionSigningSecret: "action-secret",
		slackViewOpener:          opener,
	}
	rawActionPayload := mustJSON(t, buildTestSlackActionPayload(slackActionArchive, 1001, 123456, "action-secret", now))
	req := newSignedSlackActionRequest(t, "slack-secret", rawActionPayload, now)
	rec := httptest.NewRecorder()

	svc.slackActionHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Unexpected status: got=%d body=%s", rec.Code, rec.Body.String())
	}
	if opener.calls != 1 {
		t.Fatalf("Unexpected Slack modal call count: got=%d", opener.calls)
	}
	if opener.view.CallbackID != slackActionCallbackArchive {
		t.Fatalf("Unexpected callback ID: got=%s", opener.view.CallbackID)
	}
	if opener.view.Title.Text != "Archive Finding" {
		t.Fatalf("Unexpected title: got=%s", opener.view.Title.Text)
	}
	if len(opener.view.Blocks) != 2 {
		t.Fatalf("Unexpected block count: got=%d", len(opener.view.Blocks))
	}
	if opener.view.Blocks[1].BlockID != "archive_reason" || opener.view.Blocks[1].Optional {
		t.Fatalf("Unexpected Archive reason block: got=%+v", opener.view.Blocks[1])
	}
}

func TestSlackActionRouteOpensModalWithoutCSRFToken(t *testing.T) {
	now := time.Now()
	opener := &fakeSlackViewOpener{}
	svc := &gatewayService{
		slackSigningSecret:       "slack-secret",
		slackActionSigningSecret: "action-secret",
		slackViewOpener:          opener,
	}
	rawActionPayload := mustJSON(t, buildTestSlackActionPayload(slackActionPend, 1001, 123456, "action-secret", now))
	req := newSignedSlackActionRequest(t, "slack-secret", rawActionPayload, now)
	rec := httptest.NewRecorder()

	newRouter(svc).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Unexpected status: got=%d body=%s", rec.Code, rec.Body.String())
	}
	if opener.calls != 1 {
		t.Fatalf("Unexpected Slack modal call count: got=%d", opener.calls)
	}
}

func TestSlackActionHandlerRejectsInvalidSlackSignature(t *testing.T) {
	now := time.Now()
	opener := &fakeSlackViewOpener{}
	svc := &gatewayService{
		slackSigningSecret:       "slack-secret",
		slackActionSigningSecret: "action-secret",
		slackViewOpener:          opener,
	}
	rawActionPayload := mustJSON(t, buildTestSlackActionPayload(slackActionPend, 1001, 123456, "action-secret", now))
	req := newSignedSlackActionRequest(t, "wrong-secret", rawActionPayload, now)
	rec := httptest.NewRecorder()

	svc.slackActionHandler(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("Unexpected status: got=%d body=%s", rec.Code, rec.Body.String())
	}
	if opener.calls != 0 {
		t.Fatalf("Unexpected Slack modal call count: got=%d", opener.calls)
	}
}

func TestSlackActionHandlerRejectsInvalidActionPayload(t *testing.T) {
	now := time.Now()
	opener := &fakeSlackViewOpener{}
	svc := &gatewayService{
		slackSigningSecret:       "slack-secret",
		slackActionSigningSecret: "action-secret",
		slackViewOpener:          opener,
	}
	payload := buildTestSlackActionPayload(slackActionPend, 1001, 123456, "action-secret", now)
	payload.FindingID = 999
	req := newSignedSlackActionRequest(t, "slack-secret", mustJSON(t, payload), now)
	rec := httptest.NewRecorder()

	svc.slackActionHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Unexpected status: got=%d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "invalid or expired") {
		t.Fatalf("Unexpected response body: %s", rec.Body.String())
	}
	if opener.calls != 0 {
		t.Fatalf("Unexpected Slack modal call count: got=%d", opener.calls)
	}
}

func buildTestSlackActionPayload(action string, projectID uint32, findingID uint64, secret string, now time.Time) slackActionPayload {
	payload := slackActionPayload{
		Action:    action,
		ProjectID: projectID,
		FindingID: findingID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(30 * 24 * time.Hour).Unix(),
	}
	payload.Signature = signSlackActionPayload(payload, secret)
	return payload
}

func newSignedSlackActionRequest(t *testing.T, signingSecret, rawActionPayload string, now time.Time) *http.Request {
	t.Helper()
	payload := map[string]interface{}{
		"type":       "interactive_message",
		"trigger_id": "trigger-1",
		"actions": []map[string]string{
			{
				"name":  "pend_finding",
				"value": rawActionPayload,
			},
		},
	}
	rawInteractionPayload := mustJSON(t, payload)
	body := url.Values{"payload": []string{rawInteractionPayload}}.Encode()
	req := httptest.NewRequest(http.MethodPost, slackActionEndpointPath, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, values := range signedSlackHeaders(signingSecret, []byte(body), now) {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return req
}

func signedSlackHeaders(secret string, body []byte, now time.Time) http.Header {
	timestamp := fmt.Sprintf("%d", now.Unix())
	mac := hmac.New(sha256.New, []byte(secret))
	fmt.Fprintf(mac, "%s:%s:%s", slackSignatureVersion, timestamp, string(body))
	return http.Header{
		"X-Slack-Request-Timestamp": []string{timestamp},
		"X-Slack-Signature":         []string{slackSignatureVersion + "=" + hex.EncodeToString(mac.Sum(nil))},
	}
}

func mustJSON(t *testing.T, value interface{}) string {
	t.Helper()
	buf, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	return string(buf)
}
