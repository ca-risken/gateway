package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/schema"
)

var (
	decoder                      = newDecoder()
	errRequestBodyTooLarge       = errors.New("request body too large")
	maxRequestBodyBytes    int64 = 1024 * 1024
)

func newDecoder() *schema.Decoder {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	d.ZeroEmpty(true)
	d.SetAliasTag("json")
	d.RegisterConverter([]string{}, func(input string) reflect.Value {
		return reflect.ValueOf(stringSeparator(input, ','))
	})
	return d
}

// bind bindding request parameter
func bind(out interface{}, r *http.Request) error {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		if err := bindQuery(out, r); err != nil {
			appLogger.Warnf(ctx, "Could not `bindQuery`, url=%s, err=%+v", r.URL.RequestURI(), err)
			return err
		}
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		if err := bindBodyJSON(out, r); err != nil {
			appLogger.Warnf(ctx, "Could not `bindBodyJSON`, url=%s, err=%+v", r.URL.RequestURI(), err)
			return err
		}
	default:
		appLogger.Warnf(ctx, "Unexpected HTTP Method, method=%s", r.Method)
		return fmt.Errorf("unexpected HTTP Method, method=%s", r.Method)
	}
	return nil
}

// bindQuery binding query parameter
func bindQuery(out interface{}, r *http.Request) error {
	return decoder.Decode(out, r.URL.Query())
}

// bindBodyJSON bindding body parameter binding
func bindBodyJSON(out interface{}, r *http.Request) error {
	body, err := readRequestBody(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

func readRequestBody(r *http.Request) ([]byte, error) {
	if r == nil || r.Body == nil {
		return nil, nil
	}
	limit := getMaxRequestBodyBytes()
	if r.ContentLength > limit {
		return nil, errRequestBodyTooLarge
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, limit+1))
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return nil, errRequestBodyTooLarge
		}
		return nil, err
	}
	if int64(len(body)) > limit {
		return nil, errRequestBodyTooLarge
	}
	return body, nil
}

func bufferRequestBody(r *http.Request) ([]byte, error) {
	if !requestMethodHasBody(r.Method) {
		return nil, nil
	}
	body, err := readRequestBody(r)
	if err != nil {
		return nil, err
	}
	restoreRequestBody(r, body)
	return body, nil
}

func restoreRequestBody(r *http.Request, body []byte) {
	if r == nil {
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
}

func requestMethodHasBody(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return true
	default:
		return false
	}
}

func getMaxRequestBodyBytes() int64 {
	if maxRequestBodyBytes <= 0 {
		return 1024 * 1024
	}
	return maxRequestBodyBytes
}

func stringSeparator(input string, delimiter rune) []string {
	separated := []string{}
	for _, p := range strings.Split(input, string(delimiter)) {
		if p != "" {
			separated = append(separated, p)
		}
	}
	return separated
}
