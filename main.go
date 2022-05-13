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
	Version          string
}

const (
	defaultSuffix = "Export"
	version       = "v1.1.3"
)

func main() {
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)

		return
	}

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
						ModelNamePrivate: fmt.Sprintf("%s%s", strings.ToLower(message.GetName()[:1]), message.GetName()[1:]),
						ModelName:        message.GetName(),
						PackageName:      string(file.GoPackageName),
						PathFile:         file.Desc.Path(),
						Version:          version,
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
