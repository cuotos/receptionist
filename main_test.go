package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/testutil/assert"
	"testing"
)

//TODO(dp): This is open to some serious fuzzing....
func TestExtractPorts(t *testing.T) {
	tcs := []struct {
		Title            string
		InputDockerPorts []types.Port
		Expected         []Port
	}{
		{
			"single port",
			[]types.Port{
				{PrivatePort: 80, PublicPort:  9090},
			},
			[]Port{
				{"9090", "80",""},
			},
		},
		{
			"two ports",
			[]types.Port{
				{PrivatePort: 1111, PublicPort: 2222},
				{PrivatePort: 3333, PublicPort: 4444},
			},
			[]Port{
				{"2222", "1111",""},
				{"4444", "3333",""},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Title, func(t *testing.T) {

			mockContainer := types.Container{}
			mockContainer.Ports = tc.InputDockerPorts

			actual, _ := getAllPortsFromContainer(mockContainer)

			assert.DeepEqual(t, actual, tc.Expected)
			})
	}
}

func TestCanAddNamesToPorts(t *testing.T) {

	tcs := []struct {
		Name string
		InputPorts []string
		LabelString string
		Expected []string
	}{
		{
			"single named port",
			[]string{"1111"},
			"TestPort:1111",
			[]string{
				"TestPort",
			},
		},
		{
			"two named port",
			[]string{"1111", "2222"},
			"SecondPort:2222,TestPort:1111",
			[]string{"TestPort", "SecondPort"},
		},
		// TODO: add more tests to this
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {

			for i, p := range tc.InputPorts {
				actual := getPortName(&Port{
					PrivatePort: p,
				}, tc.LabelString)
				assert.Equal(t, actual, tc.Expected[i])
			}
		})
	}
}

// There is no set ordering of the ports when they get displayed in the UI
func TestSortSliceOfPortsToBeRendered(t *testing.T) {
	tcs := []struct{
		InputPorts []*Port
		ExpectedOrder []string
	} {
		{
			[]*Port{
				&Port{"5678", "", "Second"},
				&Port{"1234", "", "First"},
			},
			[]string{"1234","5678"},
		},
	}

	for _, tc := range tcs {
		sortPorts(tc.InputPorts)

		for i, p := range tc.InputPorts {
			assert.Equal(t, p.PublicPort, tc.ExpectedOrder[i])
		}

	}

}