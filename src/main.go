package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

func main() {
	var tplPath string
	var dataPath string

	flag.StringVar(&tplPath, "template", "", "Path to the template YAML file")
	flag.StringVar(&dataPath, "data", "", "Path to the JSON configuration file")
	flag.Parse()

	if tplPath == "" || dataPath == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -template <template.yaml> -data <config.json>\n", os.Args[0])
		os.Exit(1)
	}

	// 1. Leggi il template da file
	tplBytes, err := os.ReadFile(tplPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore lettura template: %v\n", err)
		os.Exit(1)
	}

	// 2a. Costruisci la mappa di funzioni includendo Sprig e toYaml
	funcMap := sprig.TxtFuncMap()
	funcMap["toYaml"] = func(v interface{}) string {
		out, err := yaml.Marshal(v)
		if err != nil {
			return ""
		}
		return string(out)
	}

	// 2. Crea il template abilitando tutte le funzioni Sprig e toYaml
	tpl, err := template.
		New("clusterprofile").
		Funcs(funcMap).
		Parse(string(tplBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore parsing template: %v\n", err)
		os.Exit(1)
	}

	// 3. Carica i valori del cluster da cluster.json
	dataBytes, err := os.ReadFile(dataPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore lettura data file: %v\n", err)
		os.Exit(1)
	}
	var values map[string]interface{}
	if err := json.Unmarshal(dataBytes, &values); err != nil {
		fmt.Fprintf(os.Stderr, "Errore parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// 4. Esegui il rendering su stdout
	if err := tpl.Execute(os.Stdout, values); err != nil {
		fmt.Fprintf(os.Stderr, "Errore esecuzione template: %v\n", err)
		os.Exit(1)
	}
}
