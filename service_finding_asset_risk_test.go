package main

import (
	"encoding/json"
	"testing"

	"github.com/ca-risken/core/proto/finding"
)

func TestFindingAssetRiskResponseJSON(t *testing.T) {
	asset := &finding.FindingAsset{ProjectId: 1, AssetType: "VM", AssetKey: "vm", RiskInformation: []*finding.FindingAssetRisk{
		{Kind: "exploit", FindingId: 42, CveId: "CVE-2099-0001", Detail: "KEV", UpdatedAt: 100},
		{Kind: "exploit", FindingId: 42, CveId: "CVE-2099-0001", Detail: "AVAILABLE", UpdatedAt: 100},
	}}
	cases := []struct {
		name     string
		response any
	}{
		{"project", &finding.ListFindingResponse{Assets: []*finding.FindingAsset{asset}}},
		{"organization", &finding.ListFindingForOrgResponse{Assets: []*finding.FindingAsset{asset}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.response)
			if err != nil {
				t.Fatal(err)
			}
			var body struct {
				Assets []struct {
					Risks []struct {
						Kind      string `json:"kind"`
						FindingID uint64 `json:"finding_id"`
						CVEID     string `json:"cve_id"`
						Detail    string `json:"detail"`
						UpdatedAt int64  `json:"updated_at"`
					} `json:"risk_information"`
				} `json:"assets"`
			}
			if err := json.Unmarshal(data, &body); err != nil {
				t.Fatal(err)
			}
			if len(body.Assets) != 1 || len(body.Assets[0].Risks) != 2 {
				t.Fatalf("missing evidence: %s", data)
			}
			for i, risk := range body.Assets[0].Risks {
				if risk.Kind != "exploit" || risk.FindingID != 42 || risk.CVEID != "CVE-2099-0001" || risk.UpdatedAt != 100 || risk.Detail != asset.RiskInformation[i].Detail {
					t.Fatalf("unexpected evidence: %s", data)
				}
			}
		})
	}
}
