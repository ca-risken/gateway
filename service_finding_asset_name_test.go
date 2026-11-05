package main

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ca-risken/core/proto/finding"
	"github.com/stretchr/testify/require"
)

func TestAssetNameQueryBinding(t *testing.T) {
	cases := []struct {
		name, value string
		org         bool
	}{{"project partial", "stage", false}, {"organization partial", "stage", true}, {"project symbols", "25%_日本語", false}, {"organization symbols", "25%_日本語", true}}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			query := url.Values{"asset_view": {"true"}, "asset_name": {c.value}, "project_id": {"1001"}, "organization_id": {"1009"}, "limit": {"20"}}
			req := httptest.NewRequest("GET", "/?"+query.Encode(), nil)
			if c.org {
				r := &finding.ListFindingForOrgRequest{}
				require.NoError(t, bindQuery(r, req))
				require.Equal(t, c.value, r.AssetName)
				require.NoError(t, r.Validate())
			} else {
				r := &finding.ListFindingRequest{}
				require.NoError(t, bindQuery(r, req))
				require.Equal(t, c.value, r.AssetName)
				require.NoError(t, r.Validate())
			}
		})
	}
}
