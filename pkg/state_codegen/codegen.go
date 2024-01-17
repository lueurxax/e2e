package main

import (
	"bytes"
	"embed"
	"flag"
	"go/format"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"

	"github.com/lueurxax/e2e/pkg/state_codegen/models"
)

var (
	configPath = flag.String("config", "./.stresscfg.yaml", "path to config file")
)

//go:embed selector.gotpl stated_method_wrapper.gotpl state.gotpl stress_storage.gotpl
var f embed.FS

func main() {
	cfg := struct {
		Parameters models.Params
	}{}
	flag.Parse()

	data, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	tmpl, err := template.New("").ParseFS(
		f,
		"selector.gotpl",
		"stated_method_wrapper.gotpl",
		"state.gotpl",
		"stress_storage.gotpl")
	if err != nil {
		panic(err)
	}

	if err := execute(tmpl, cfg.Parameters.PathToGenerated, "selector", cfg.Parameters); err != nil {
		panic(err)
	}

	if err := execute(tmpl, cfg.Parameters.PathToGenerated, "stated_method_wrapper", cfg.Parameters); err != nil {
		panic(err)
	}

	if err := execute(tmpl, cfg.Parameters.PathToGenerated, "state", cfg.Parameters); err != nil {
		panic(err)
	}

	if err := execute(tmpl, cfg.Parameters.PathToGenerated, "stress_storage", cfg.Parameters); err != nil {
		panic(err)
	}
}

func execute(tmpl *template.Template, path, name string, params models.Params) error {
	buf := bytes.NewBuffer(make([]byte, 0))
	err := tmpl.ExecuteTemplate(buf, name+".gotpl", params)
	if err != nil {
		return err
	}

	res, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile(path+name+".go", res, 0600)
}
