package main

import (
	"bytes"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Make sure the index page template can render, provided it as a model of data.
// This is only testing the template renders not its contents.
func TestTemplateRendering(t *testing.T) {
	model := Model{
		Containers: []Container{
			{
				Ports: []*Port{
					{9999, 8888, "TestPort", "/"},
					{9999, 8888, "", "/"},
				},
				Name:  "TestContainer",
				Image: "test/dockerimage",
			},
		},
	}

	rendered := &bytes.Buffer{}

	err := Tpl.Execute(rendered, model)
	assert.NoError(t, err)

	// HTML is XML, make sure it can be parsed
	// TODO: check if there is a stdlib xhml parser
	err = xml.Unmarshal(rendered.Bytes(), new(interface{}))
	assert.NoError(t, err)
}