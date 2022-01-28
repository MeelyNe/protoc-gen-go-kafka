package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

//go:embed handler.kafka.tmpl
var handlerKafkaTmpl string

type gen struct {
	ModelNamePrivate string
	ModelName        string
	PackageName      string
	PathFile         string
}

const defaultSuffix = "Export"

func main() {
	var flags flag.FlagSet
	suffix := flags.String("suffix", defaultSuffix, "")
	protoc := protogen.Options{
		ParamFunc: flags.Set,
	}
	protoc.Run(func(plugin *protogen.Plugin) error {
		if *suffix == "" {
			*suffix = defaultSuffix
		}

		for _, file := range plugin.Files {
			for _, message := range file.Proto.GetMessageType() {
				if strings.HasSuffix(message.GetName(), *suffix) {
					tmpl, err := parseTemplates(&gen{
						ModelNamePrivate: strings.ToLower(message.GetName()),
						ModelName:        message.GetName(),
						PackageName:      string(file.GoPackageName),
						PathFile:         file.Desc.Path(),
					})
					if err != nil {
						plugin.Error(err)

						continue
					}
					msgName := strings.ToLower(strings.Replace(message.GetName(), *suffix, "", 2))
					filename := fmt.Sprintf("%s_%s.kafka.go", file.GeneratedFilenamePrefix, msgName)

					genFile := plugin.NewGeneratedFile(filename, file.GoImportPath)

					if _, err = genFile.Write([]byte(tmpl)); err != nil {
						plugin.Error(err)

						continue
					}
				}
			}
		}

		return nil
	})
}

func parseTemplates(tmplData interface{}) (str string, err error) {
	tmpl, err := template.New("").Parse(handlerKafkaTmpl)
	if err != nil {
		return
	}

	var content bytes.Buffer

	err = tmpl.Execute(&content, tmplData)
	if err != nil {
		return
	}

	return content.String(), nil
}
