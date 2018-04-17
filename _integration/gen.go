// +build none

package main

import (
	"html/template"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("generator", "Generate mock integration testing data.")

	services := app.Command("services", "Generate mock service data.")
	count := services.Flag("count", "Count of entries to create.").Short('c').Default("500").Int()

	rand.Seed(time.Now().UnixNano())
	args := os.Args[1:]
	switch kingpin.MustParse(app.Parse(args)) {
	case services.FullCommand():
		genServices(*count)
	default:
		app.Usage(args)
		os.Exit(2)
	}
}

const serviceTmpl = `
# autogenerated: do not edit!
# source:{{ range .args }} {{ . }}{{end}}
{{ range .names -}}
apiVersion: v1
kind: Service
metadata:
  name: {{ . }}
spec:
  ports:
  - port: 80
    protocol: TCP
---
{{ end}}
`

func genServices(n int) {
	var names []string
	for i := 0; i < n; i++ {
		name := namesgenerator.GetRandomName(0)
		name = strings.Replace(name, "_", "-", -1) // must be a valid rfc 1035 value
		names = append(names, name)
	}

	t, err := template.New("services").Parse(serviceTmpl[1:])
	check(err)
	check(t.Execute(os.Stdout, map[string]interface{}{
		"names": names,
		"args":  os.Args,
	}))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}