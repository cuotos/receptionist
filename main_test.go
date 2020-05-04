package main

import (
	"github.com/docker/docker/pkg/testutil/assert"
	"testing"
)

func TestExtractPorts(t *testing.T) {
	tcs := []struct{
		Title string
		Input string
		Expected []Port
	}{
		{
			"single port",
			"9090,",
			[]Port{
				{"9090",""},
			},
		},
		{
			"two ports",
			"9090,10101",
			[]Port{
				{"9090", ""},
				{"10101", ""},
			},
		},
		{
			"single port w/ trailing comma",
			"9090,",
			[]Port{
				{"9090", ""},
			},
		},
		{
			"single port w/ leading and trailing comma",
			",9090,",
			[]Port{
				{"9090", ""},
			},
		},
		{
			"single named port",
			"UI:9090",
			[]Port{
				{"9090", "UI"},
			},
		},
		{
			"multiple named ports",
			"UI:9090,API:10101",
			[]Port{
				{"9090", "UI"},
				{ "10101", "API"},
			},
		},
		{
			"multiple named ports w/ missing name",
			"UI:9090,API:10101,:11111",
			[]Port{
				{"9090", "UI"},
				{ "10101", "API"},
				{ "11111", ""},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Title, func(t *testing.T) {
			actual, _ := extractPorts(tc.Input)
			assert.DeepEqual(t, actual, tc.Expected)
		})

	}
}