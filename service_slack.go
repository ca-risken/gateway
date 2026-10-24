package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	slackActionEndpointPath = "/api/v1/slack/actions"
	slackActionPend         = "pend"
	slackActionArchive      = "archive"

	slackActionCallbackPend    = "risken_slack_action_pend"
	slackActionCallbackArchive = "risken_slack_action_archive"

	slackSignatureVersion  = "v0"
	slackSignatureMaxAge   = 5 * time.Minute
	slackHTTPClientTimeout       = 10 * time.Second
	slackViewsOpenURL            = "https://slack.com/api/views.open"
	slackAPIResponseMaxBytes     = 64 * 1024
	slackAPIResponseLogMaxBytes  = 256
)

var (
	errSlackSigningSecretNotConfigured       = errors.New("slack signing secret is not configured")
	errSlackActionSigningSecretNotConfigured = errors.New("slack action signing secret is not configured")
	errSlackBotTokenNotConfigured            = errors.New("slack bot token is not configured")
	errSlackRequestTimestampMissing          = errors.New("slack request timestamp is missing")
	errSlackRequestTimestampExpired          = errors.New("slack request timestamp is expired")
	errSlackRequestSignatureInvalid          = errors.New("slack request signature is invalid")
	errSlackInteractionPayloadMissing        = errors.New("slack interaction payload is missing")
	errSlackActionPayloadMissing             = errors.New("slack action payload is missing")
	errSlackActionPayloadInvalid             = errors.New("slack action payload is invalid")
	errSlackActionPayloadExpired             = errors.New("slack action payload is expired")
	errSlackActionPayloadSignatureInvalid    = errors.New("slack action payload signature is invalid")
)

type slackViewOpener interface {
	OpenView(ctx context.Context, triggerID string, view slackModalView) error
}

type slackAPIClient struct {
	token      string
	httpClient *http.Client
}

type slackOpenViewRequest struct {
	TriggerID string         `json:"trigger_id"`
	View      slackModalView `json:"view"`
}

type slackAPIResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type slackInteractionPayload struct {
	Type      string                   `json:"type"`
	TriggerID string                   `json:"trigger_id"`
	Actions   []slackInteractionAction `json:"actions"`
}

type slackInteractionAction struct {
	Name     string `json:"name"`
	ActionID string `json:"action_id"`
	Value    string `json:"value"`
}

type slackActionPayload struct {
	Action    string `json:"action"`
	ProjectID uint32 `json:"project_id"`
	FindingID uint64 `json:"finding_id"`
	IssuedAt  int64  `json:"issued_at"`
	ExpiresAt int64  `json:"expires_at"`
	Signature string `json:"signature"`
}

type slackModalView struct {
	Type            string             `json:"type"`
	CallbackID      string             `json:"callback_id"`
	PrivateMetadata string             `json:"private_metadata"`
	Title           slackTextObject    `json:"title"`
	Submit          slackTextObject    `json:"submit"`
	Close           slackTextObject    `json:"close"`
	Blocks          []slackBlockObject `json:"blocks"`
}

type slackBlockObject struct {
	Type     string              `json:"type"`
	BlockID  string              `json:"block_id,omitempty"`
	Text     *slackTextObject    `json:"text,omitempty"`
	Label    *slackTextObject    `json:"label,omitempty"`
	Element  *slackElementObject `json:"element,omitempty"`
	Optional bool                `json:"optional,omitempty"`
}

type slackElementObject struct {
	Type        string           `json:"type"`
	ActionID    string           `json:"action_id"`
	Placeholder *slackTextObject `json:"placeholder,omitempty"`
	Multiline   bool             `json:"multiline,omitempty"`
}

type slackTextObject struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func readSlackAPIResponseBody(r io.Reader) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(r, slackAPIResponseMaxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > slackAPIResponseMaxBytes {
		return nil, fmt.Errorf("slack api response body exceeds max size %d", slackAPIResponseMaxBytes)
	}
	return body, nil
}

func truncateSlackResponseForLog(body []byte) string {
	if len(body) <= slackAPIResponseLogMaxBytes {
		return string(body)
	}
	return string(body[:slackAPIResponseLogMaxBytes]) + "...(truncated)"
}

func newSlackAPIClient(token string) *slackAPIClient {
	return &slackAPIClient{
		token:      token,
		httpClient: &http.Client{Timeout: slackHTTPClientTimeout},
	}
}

func (c *slackAPIClient) OpenView(ctx context.Context, triggerID string, view slackModalView) error {
	if c.token == "" {
		return errSlackBotTokenNotConfigured
	}
	body, err := json.Marshal(slackOpenViewRequest{
		TriggerID: triggerID,
		View:      view,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, slackViewsOpenURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := readSlackAPIResponseBody(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack views.open status=%d body=%s", resp.StatusCode, truncateSlackResponseForLog(respBody))
	}
	var slackResp slackAPIResponse
	if err := json.Unmarshal(respBody, &slackResp); err != nil {
		return err
	}
	if !slackResp.OK {
		return fmt.Errorf("slack views.open error=%s", slackResp.Error)
	}
	return nil
}

func (g *gatewayService) slackActionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		appLogger.Errorf(ctx, "Failed to read Slack action request body, err=%+v", err)
		http.Error(w, "Could not read body", http.StatusInternalServerError)
		return
	}

	if err := verifySlackRequestSignature(g.slackSigningSecret, r.Header, body, time.Now()); err != nil {
		appLogger.Warnf(ctx, "Rejected Slack action request signature, err=%+v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	interaction, rawActionPayload, err := parseSlackInteractionPayload(body)
	if err != nil {
		appLogger.Warnf(ctx, "Failed to parse Slack interaction payload, err=%+v", err)
		writeSlackEphemeralResponse(w, "Failed to parse Slack action payload.")
		return
	}

	actionPayload, err := parseAndVerifySlackActionPayload(rawActionPayload, g.slackActionSigningSecret, time.Now())
	if err != nil {
		appLogger.Warnf(ctx, "Rejected Slack action payload, err=%+v", err)
		writeSlackEphemeralResponse(w, "This Slack action is invalid or expired. Please open RISKEN and confirm the Finding.")
		return
	}

	view, err := buildSlackActionModalView(actionPayload, rawActionPayload)
	if err != nil {
		appLogger.Warnf(ctx, "Failed to build Slack action modal view, err=%+v", err)
		writeSlackEphemeralResponse(w, "Failed to open the Slack action form.")
		return
	}
	if g.slackViewOpener == nil {
		appLogger.Errorf(ctx, "Slack view opener is not configured")
		writeSlackEphemeralResponse(w, "Slack action is not configured.")
		return
	}
	if err := g.slackViewOpener.OpenView(ctx, interaction.TriggerID, view); err != nil {
		appLogger.Errorf(ctx, "Failed to open Slack modal, err=%+v", err)
		writeSlackEphemeralResponse(w, "Failed to open the Slack action form.")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func verifySlackRequestSignature(secret string, header http.Header, body []byte, now time.Time) error {
	if secret == "" {
		return errSlackSigningSecretNotConfigured
	}
	timestamp := header.Get("X-Slack-Request-Timestamp")
	if timestamp == "" {
		return errSlackRequestTimestampMissing
	}
	requestUnix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return err
	}
	requestTime := time.Unix(requestUnix, 0)
	if now.Sub(requestTime) > slackSignatureMaxAge || requestTime.Sub(now) > slackSignatureMaxAge {
		return errSlackRequestTimestampExpired
	}

	mac := hmac.New(sha256.New, []byte(secret))
	fmt.Fprintf(mac, "%s:%s:%s", slackSignatureVersion, timestamp, string(body))
	want := slackSignatureVersion + "=" + hex.EncodeToString(mac.Sum(nil))
	got := header.Get("X-Slack-Signature")
	if !hmac.Equal([]byte(got), []byte(want)) {
		return errSlackRequestSignatureInvalid
	}
	return nil
}

func parseSlackInteractionPayload(body []byte) (*slackInteractionPayload, string, error) {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, "", err
	}
	rawPayload := values.Get("payload")
	if rawPayload == "" {
		return nil, "", errSlackInteractionPayloadMissing
	}
	var payload slackInteractionPayload
	if err := json.Unmarshal([]byte(rawPayload), &payload); err != nil {
		return nil, "", err
	}
	if payload.TriggerID == "" {
		return nil, "", errSlackInteractionPayloadMissing
	}
	if len(payload.Actions) == 0 || payload.Actions[0].Value == "" {
		return nil, "", errSlackActionPayloadMissing
	}
	return &payload, payload.Actions[0].Value, nil
}

func parseAndVerifySlackActionPayload(rawPayload, secret string, now time.Time) (*slackActionPayload, error) {
	if secret == "" {
		return nil, errSlackActionSigningSecretNotConfigured
	}
	var payload slackActionPayload
	if err := json.Unmarshal([]byte(rawPayload), &payload); err != nil {
		return nil, err
	}
	if payload.Action != slackActionPend && payload.Action != slackActionArchive {
		return nil, errSlackActionPayloadInvalid
	}
	if payload.ProjectID == 0 || payload.FindingID == 0 || payload.IssuedAt == 0 || payload.ExpiresAt == 0 || payload.Signature == "" {
		return nil, errSlackActionPayloadInvalid
	}
	if now.Unix() > payload.ExpiresAt {
		return nil, errSlackActionPayloadExpired
	}
	want := signSlackActionPayload(payload, secret)
	if !hmac.Equal([]byte(payload.Signature), []byte(want)) {
		return nil, errSlackActionPayloadSignatureInvalid
	}
	return &payload, nil
}

func signSlackActionPayload(payload slackActionPayload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	fmt.Fprintf(mac, "%s:%d:%d:%d:%d", payload.Action, payload.ProjectID, payload.FindingID, payload.IssuedAt, payload.ExpiresAt)
	return hex.EncodeToString(mac.Sum(nil))
}

func buildSlackActionModalView(payload *slackActionPayload, rawPayload string) (slackModalView, error) {
	switch payload.Action {
	case slackActionPend:
		return buildPendSlackModalView(payload, rawPayload), nil
	case slackActionArchive:
		return buildArchiveSlackModalView(payload, rawPayload), nil
	default:
		return slackModalView{}, errSlackActionPayloadInvalid
	}
}

func buildPendSlackModalView(payload *slackActionPayload, rawPayload string) slackModalView {
	return slackModalView{
		Type:            "modal",
		CallbackID:      slackActionCallbackPend,
		PrivateMetadata: rawPayload,
		Title:           plainText("Pending Finding"),
		Submit:          plainText("Submit"),
		Close:           plainText("Cancel"),
		Blocks: []slackBlockObject{
			contextSection(payload),
			{
				Type:    "input",
				BlockID: "pend_deadline",
				Label:   &slackTextObject{Type: "plain_text", Text: "PEND deadline"},
				Element: &slackElementObject{
					Type:        "datepicker",
					ActionID:    "date",
					Placeholder: &slackTextObject{Type: "plain_text", Text: "Select a deadline"},
				},
			},
			{
				Type:     "input",
				BlockID:  "pend_note",
				Label:    &slackTextObject{Type: "plain_text", Text: "Note"},
				Optional: true,
				Element: &slackElementObject{
					Type:        "plain_text_input",
					ActionID:    "text",
					Placeholder: &slackTextObject{Type: "plain_text", Text: "Add a note"},
					Multiline:   true,
				},
			},
		},
	}
}

func buildArchiveSlackModalView(payload *slackActionPayload, rawPayload string) slackModalView {
	return slackModalView{
		Type:            "modal",
		CallbackID:      slackActionCallbackArchive,
		PrivateMetadata: rawPayload,
		Title:           plainText("Archive Finding"),
		Submit:          plainText("Submit"),
		Close:           plainText("Cancel"),
		Blocks: []slackBlockObject{
			contextSection(payload),
			{
				Type:    "input",
				BlockID: "archive_reason",
				Label:   &slackTextObject{Type: "plain_text", Text: "Archive reason"},
				Element: &slackElementObject{
					Type:        "plain_text_input",
					ActionID:    "text",
					Placeholder: &slackTextObject{Type: "plain_text", Text: "Describe why this Finding can be archived"},
					Multiline:   true,
				},
			},
		},
	}
}

func contextSection(payload *slackActionPayload) slackBlockObject {
	return slackBlockObject{
		Type: "section",
		Text: &slackTextObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Project ID:* %d\n*Finding ID:* %d", payload.ProjectID, payload.FindingID),
		},
	}
}

func plainText(text string) slackTextObject {
	return slackTextObject{
		Type: "plain_text",
		Text: text,
	}
}

func writeSlackEphemeralResponse(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"response_type": "ephemeral",
		"text":          text,
	})
}
