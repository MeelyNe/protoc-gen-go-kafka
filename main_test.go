package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func Test_main(t *testing.T) {
	var suffix = "Out"
	ass := assert.New(t)
	binary, err := os.ReadFile("./test/stdin")
	ass.NoError(err)
	req := &pluginpb.CodeGeneratorRequest{}
	err = proto.Unmarshal(binary, req)
	ass.NoError(err)

	plugin, err := protogen.Options{}.New(req)
	ass.NoError(err)
	err = func(plugin *protogen.Plugin, t *testing.T) error {
		for _, file := range plugin.Files {
			for _, message := range file.Proto.GetMessageType() {
				ass.Equal(message.GetName(), "OutgoingMessageOut")
				if strings.HasSuffix(message.GetName(), suffix) {
					ass.Equal(&gen{
						ModelNamePrivate: strings.ToLower(message.GetName()),
						ModelName:        message.GetName(),
						PackageName:      string(file.GoPackageName),
						PathFile:         file.Desc.Path(),
					}, &gen{
						ModelNamePrivate: "outgoingmessageout",
						ModelName:        "OutgoingMessageOut",
						PackageName:      "outgoingv1",
						PathFile:         "outgoing_topic/v1/outgoing.proto",
					})
					tmpl, err := parseTemplates(&gen{
						ModelNamePrivate: strings.ToLower(message.GetName()),
						ModelName:        message.GetName(),
						PackageName:      string(file.GoPackageName),
						PathFile:         file.Desc.Path(),
					})
					if err != nil {
						ass.NoError(err)
						plugin.Error(err)

						continue
					}
					msgName := strings.ToLower(strings.Replace(message.GetName(), suffix, "", 2))
					filename := fmt.Sprintf("%s_%s.kafka.go", file.GeneratedFilenamePrefix, msgName)

					genFile := plugin.NewGeneratedFile(filename, file.GoImportPath)

					if _, err = genFile.Write([]byte(tmpl)); err != nil {
						ass.NoError(err)
						plugin.Error(err)

						continue
					}
				}
			}
		}

		return nil
	}(plugin, t)

	ass.NoError(err)

}

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
