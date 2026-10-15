package main

import (
	"encoding/json"
	"github.com/ca-risken/core/proto/finding"
	"net/http/httptest"
	"testing"
)

func TestFindingAssetQueryBinding(t *testing.T) {
	cases := []struct {
		name  string
		org   bool
		query string
		view  bool
	}{
		{"project groups", false, "project_id=1001&asset_view=true&from_score=0.8&status=1&sort=finding_count", true},
		{"organization groups", true, "organization_id=1009&asset_view=true&from_score=0.8&status=1&sort=finding_count", true},
		{"project expansion", false, "project_id=1001&asset_view=false&asset_type=Workload&asset_key=%5B%22cluster%22%2C%22ns%22%2C%22Deployment%22%2C%22app%22%5D&asset_project_id=1001", false},
		{"organization expansion", true, "organization_id=1009&asset_view=false&asset_type=Workload&asset_key=%5B%22cluster%22%2C%22ns%22%2C%22Deployment%22%2C%22app%22%5D&asset_project_id=1001", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/?"+c.query, nil)
			var view bool
			var key, typ string
			var project uint32
			if c.org {
				req := &finding.ListFindingForOrgRequest{}
				if err := bindQuery(req, r); err != nil {
					t.Fatal(err)
				}
				view, key, typ, project = req.AssetView, req.AssetKey, req.AssetType, req.AssetProjectId
			} else {
				req := &finding.ListFindingRequest{}
				if err := bindQuery(req, r); err != nil {
					t.Fatal(err)
				}
				view, key, typ, project = req.AssetView, req.AssetKey, req.AssetType, req.AssetProjectId
			}
			if view != c.view {
				t.Fatalf("asset_view=%v", view)
			}
			if !c.view && (key != `["cluster","ns","Deployment","app"]` || typ != "Workload" || project != 1001) {
				t.Fatalf("expansion filters lost: %s %s %d", typ, key, project)
			}
		})
	}
}
func TestFindingAssetResponseJSON(t *testing.T) {
	cases := []struct {
		name     string
		response any
	}{
		{"project", &finding.ListFindingResponse{Total: 1, Assets: []*finding.FindingAsset{{ProjectId: 1001, AssetType: "VM", AssetKey: "instance", AssetName: "vm", MaxScore: 0.9, FindingCount: 3}}}},
		{"organization", &finding.ListFindingForOrgResponse{Total: 1, Assets: []*finding.FindingAsset{{ProjectId: 1001, AssetType: "VM", AssetKey: "instance", AssetName: "vm", MaxScore: 0.9, FindingCount: 3}}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, err := json.Marshal(c.response)
			if err != nil {
				t.Fatal(err)
			}
			var body struct {
				Total  int `json:"total"`
				Assets []struct {
					Key     string `json:"asset_key"`
					Count   int    `json:"finding_count"`
					Project int    `json:"project_id"`
				} `json:"assets"`
			}
			if err := json.Unmarshal(b, &body); err != nil {
				t.Fatal(err)
			}
			if body.Total != 1 || len(body.Assets) != 1 || body.Assets[0].Key != "instance" || body.Assets[0].Count != 3 || body.Assets[0].Project != 1001 {
				t.Fatalf("unexpected response: %s", b)
			}
		})
	}
}
