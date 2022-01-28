package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseTemplates(t *testing.T) {
	var crateFileTable = []struct {
		label       string
		tmplBody    string
		tmplData    *gen
		wantContent string
	}{
		{
			tmplBody:    `templates GEN-TEST-1 {{.ModelNamePrivate}} {{.ModelName}} {{.PackageName}} {{.PathFile}}`,
			tmplData:    &gen{ModelNamePrivate: "test1", ModelName: "Test1", PackageName: "testv1", PathFile: "/to.proto"},
			wantContent: `templates GEN-TEST-1 test1 Test1 testv1 /to.proto`,
		},
		{
			tmplBody:    `templates GEN-TEST-2 {{.ModelNamePrivate}} {{.ModelName}} {{.PackageName}} {{.PathFile}}`,
			tmplData:    &gen{ModelNamePrivate: "abc", ModelName: "ABC", PackageName: "abcv1", PathFile: "/abc.proto"},
			wantContent: `templates GEN-TEST-2 abc ABC abcv1 /abc.proto`,
		},
	}

	ass := assert.New(t)

	for _, tt := range crateFileTable {
		t.Run(tt.label, func(t *testing.T) {
			handlerKafkaTmpl = tt.tmplBody
			result, err := parseTemplates(tt.tmplData)
			ass.NoError(err)

			ass.Equal(result, tt.wantContent)
		})
	}
}
