package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/schema"
)

var (
	decoder = newDecoder()
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
	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	return err
	// }
	// return json.Unmarshal(body, out) // Need to read the body several times.(json.Decoder is only once)
	return json.NewDecoder(r.Body).Decode(out)
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
