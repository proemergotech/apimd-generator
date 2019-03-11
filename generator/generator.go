package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(d Definitons) {
	c := newCollector()
	doc := c.collect(d)

	t, err := template.
		New("").
		Funcs(template.FuncMap{
			"dict": func(v ...interface{}) map[string]interface{} {
				result := make(map[string]interface{}, len(v)/2)
				for i := 0; i < len(v)-1; i += 2 {
					result[fmt.Sprintf("%s", v[i])] = v[i+1]
				}

				return result
			},
			"isValue": func(v interface{}) bool {
				_, ok := v.(*DocValue)
				return ok
			},
			"isArray": func(v interface{}) bool {
				_, ok := v.([]interface{})
				return ok
			},
			"add": func(i1 int, i2 int) int {
				return i1 + i2
			},
			"indent": func(i int) string {
				return strings.Repeat(" ", i)
			},
			"dig3": func(i int) string {
				result := strconv.Itoa(i)
				return strings.Repeat("0", 3-len(result)) + result
			},
		}).
		Parse(apimdTmpl)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	buf := &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "base", doc)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	apimd, err := filepath.Abs(d.OutputPath())
	if err != nil {
		log.Fatalf("%+v", err)
	}

	err = ioutil.WriteFile(apimd, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	log.Print("Updated " + apimd + ":1")
}
