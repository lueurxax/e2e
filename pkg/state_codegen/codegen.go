package main

import (
	"bytes"
	"flag"
	"go/format"
	"io/ioutil"
	"text/template"

	"gopkg.in/yaml.v3"

	"git.proksy.io/golang/e2e/pkg/state_codegen/models"
)

var (
	configPath = flag.String("config", "./.stresscfg.yaml", "path to config file")
)

func main() {
	cfg := struct {
		Parameters models.Params
	}{}

	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	tmpl, err := template.New("").ParseFiles(
		"pkg/state_codegen/selector.gotpl",
		"pkg/state_codegen/stated_method_wrapper.gotpl",
		"pkg/state_codegen/state.gotpl",
		"pkg/state_codegen/stress_storage.gotpl")
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

	return ioutil.WriteFile(path+name+".go", res, 0600)
}
