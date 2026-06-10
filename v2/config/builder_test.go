package config

import (
	"context"
	"strings"
	"testing"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func TestBuildConfigPreservesInputGroupsAndRules(t *testing.T) {
	input := &option.Options{
		Outbounds: []option.Outbound{
			{
				Type: C.TypeSOCKS,
				Tag:  "node-a",
				Options: &option.SOCKSOutboundOptions{
					ServerOptions: option.ServerOptions{
						Server:     "127.0.0.1",
						ServerPort: 1080,
					},
				},
			},
			{
				Type: C.TypeSelector,
				Tag:  "MAIN",
				Options: &option.SelectorOutboundOptions{
					Outbounds: []string{"MANUAL", "direct"},
					Default:   "MANUAL",
				},
			},
			{
				Type: C.TypeSelector,
				Tag:  "MANUAL",
				Options: &option.SelectorOutboundOptions{
					Outbounds: []string{"node-a"},
					Default:   "node-a",
				},
			},
		},
		Route: &option.RouteOptions{
			Rules: []option.Rule{
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						RawDefaultRule: option.RawDefaultRule{
							DomainSuffix: []string{"telegram.org"},
						},
						RuleAction: option.RuleAction{
							Action: C.RuleActionTypeRoute,
							RouteOptions: option.RouteActionOptions{
								Outbound: "MAIN",
							},
						},
					},
				},
			},
			Final: "MAIN",
		},
	}

	result, err := BuildConfig(context.Background(), DefaultHiddifyOptions(), &ReadOptions{Options: input})
	if err != nil {
		t.Fatalf("BuildConfig failed: %v", err)
	}

	content, err := result.MarshalJSONContext(context.Background())
	if err != nil {
		t.Fatalf("MarshalJSONContext failed: %v", err)
	}

	json := string(content)
	for _, expected := range []string{
		`"tag": "MAIN"`,
		`"tag": "MANUAL"`,
		`"outbound": "MAIN"`,
		`"domain_suffix": "telegram.org"`,
		`"final": "MAIN"`,
	} {
		if !strings.Contains(json, expected) {
			t.Fatalf("expected built config to contain %s, got: %s", expected, json)
		}
	}
}