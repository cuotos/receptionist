package main

import (
	"github.com/docker/docker/pkg/testutil/assert"
	"testing"
)

func TestExtractPorts(t *testing.T) {
	tcs := []struct{
		Title     string
		Input     string
		ExpectAll bool
		Expected  []Port
	}{
		{
			"single port",
			"9090,",
			false,
			[]Port{
				{"9090",""},
			},
		},
		{
			"two ports",
			"9090,10101",
			false,
			[]Port{
				{"9090", ""},
				{"10101", ""},
			},
		},
		{
			"single port w/ trailing comma",
			"9090,",
			false,
			[]Port{
				{"9090", ""},
			},
		},
		{
			"single port w/ leading and trailing comma",
			",9090,",
			false,
			[]Port{
				{"9090", ""},
			},
		},
		{
			"single named port",
			"UI:9090",
			false,
			[]Port{
				{"9090", "UI"},
			},
		},
		{
			"multiple named ports",
			"UI:9090,API:10101",
			false,
			[]Port{
				{"9090", "UI"},
				{ "10101", "API"},
			},
		},
		{
			"multiple named ports w/ missing name",
			"UI:9090,API:10101,:11111",
			false,
			[]Port{
				{"9090", "UI"},
				{ "10101", "API"},
				{ "11111", ""},
			},
		},
		{
			"all ports should be exposed",
			"ALL",
			true,
			nil,
		},
		{
			"named port plus all port flag",
			"API:1111,ALL",
			true,
			[]Port{
				{"1111","API"},
			},
		},
		{
			"no name port and all port flag",
			":1111,ALL",
			true,
			[]Port{
				{"1111", ""},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Title, func(t *testing.T) {
			actual, all, _ := parsePortsLabel(tc.Input)
			assert.DeepEqual(t, actual, tc.Expected)
			assert.Equal(t, all, tc.ExpectAll)
		})
	}
}

//func TestAllPortsFlag(t *testing.T) {
//	tcs := []struct{
//		InputContainer types.ContainerJSON
//		ExpectedPorts []string
//	}{
//		{
//			types.ContainerJSON{
//				Config: &container.Config{
//					Labels: map[string]string{"RECEP": "test"},
//				},
//			},
//			[]string{""},
//		},
//	}
//
//	for _, tc := range tcs {
//		actual :=
//	}
//}