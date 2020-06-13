package main

import (
	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

//TODO(dp): This is open to some serious fuzzing....

// Extract the
func TestExtractPorts(t *testing.T) {
	tcs := []struct {
		Title            string
		InputDockerPorts []types.Port
		Expected         []*Port
	}{
		{
			"single port",
			[]types.Port{
				{PrivatePort: 80, PublicPort: 9090},
			},
			[]*Port{
				&Port{9090, 80, "", "/"},
			},
		},
		{
			"two ports",
			[]types.Port{
				{PrivatePort: 1111, PublicPort: 2222},
				{PrivatePort: 3333, PublicPort: 4444},
			},
			[]*Port{
				&Port{2222, 1111, "", "/"},
				&Port{4444, 3333, "", "/"},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Title, func(t *testing.T) {

			mockContainer := types.Container{}
			mockContainer.Labels = map[string]string{"RECEPTIONIST": ""}
			mockContainer.Ports = tc.InputDockerPorts

			actual, err := getAllWantedPortsFromContainer(mockContainer)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, actual, tc.Expected)
		})
	}
}

// If the RECPTIONIST label contains a <string>:<port> then assign that name to the port for use in the UI.
func TestCanPopulatePortWithNameAndPath(t *testing.T) {

	tcs := []struct {
		InputPorts   []uint16
		LabelString string
		Expected    []Port
	}{
		{
			[]uint16{1111},
			"TestPort:1111",
			[]Port{
				{
					Name: "TestPort",
					Path: "/",
				},
			},
		},
		{
			[]uint16{1111},
			"TestPort:1111:/thePath",
			[]Port{
				{
					Name: "TestPort",
					Path: "/thePath",
				},
			},
		},
		{
			[]uint16{1111},
			"TestPort:1111:/thePath",
			[]Port{
				{
					Name: "TestPort",
					Path: "/thePath",
				},
			},
		},
		{
			[]uint16{1111,2222},
			"TestPort:2222",
			[]Port{
				{},
				{
					Name: "TestPort",
					Path: "/",
				},
			},
		},
		{
			[]uint16{1111},
			"TestPort:1111:path",
			[]Port{{Name: "TestPort", Path: "/path"}},
		},
	}

	for _, tc := range tcs {
		for i, p := range tc.InputPorts {
			p := &Port{PrivatePort: p}

			err := populatePortMetaData(p, tc.LabelString)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			assert.Equal(t, p.Name, tc.Expected[i].Name)
			assert.Equal(t, p.Path, tc.Expected[i].Path)
		}
	}
}

func TestParseLabel(t *testing.T) {
	tcs := []struct{
		LabelString string
		Expected    []LabelElement
	}{
		{
			"TestPort:1111",
			[]LabelElement{
				{
					"TestPort",
					"1111",
					"/",
				},
			},
		},
		{
			"TestPort:1111,SecondPort:2222",
			[]LabelElement{
				{
					"TestPort",
					"1111",
					"/",
				},
				{
					"SecondPort",
					"2222",
					"/",
				},
			},
		},
		{
			"TestPort:1111:/aPath",
			[]LabelElement{
				{
					"TestPort",
					"1111",
					"/aPath",
				},
			},
		},
		{
			"TestPort:1111:/thePath,SecondPort:2222:/aPath",
			[]LabelElement{
				{
					"TestPort",
					"1111",
					"/thePath",
				},
				{
					"SecondPort",
					"2222",
					"/aPath",
				},
			},
		},
		{
			"TestPort:1111:/a/more/complex/path",
			[]LabelElement{
				{
					"TestPort",
					"1111",
					"/a/more/complex/path",
				},
			},
		},
		{
			"TestPort:1111:/a/more/complex/path,AnotherPort:2222",
			[]LabelElement{
				{
					"TestPort",
					"1111",
					"/a/more/complex/path",
				},
				{
					"AnotherPort",
					"2222",
					"/",
				},
			},
		},
		{
			":1111:/a/more/complex/path",
			[]LabelElement{
				{
					"",
					"1111",
					"/a/more/complex/path",
				},
			},
		},
		{
			":1111:",
			[]LabelElement{
				{
					"",
					"1111",
					"/",
				},
			},
		},
		{
			"",
			[]LabelElement{},
		},
	}

	for _, tc := range tcs {
		actual, _ := extractElementsFromLabel(tc.LabelString)
		for i, e := range tc.Expected {
			assert.Equal(t, actual[i], e)
		}
	}
}

// There is no set ordering of the ports when they get displayed in the UI
func TestSortSliceOfPortsToBeRendered(t *testing.T) {
	tcs := []struct {
		InputPorts    []*Port
		ExpectedOrder []uint16
	}{
		{
			[]*Port{
				&Port{5678, 0, "Second", "/"},
				&Port{1234, 0, "First", "/"},
			},
			[]uint16{1234, 5678},
		},
	}

	for _, tc := range tcs {
		sortPorts(tc.InputPorts)

		for i, p := range tc.InputPorts {
			assert.Equal(t, p.PublicPort, tc.ExpectedOrder[i])
		}
	}
}
