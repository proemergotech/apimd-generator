package main

import (
	"io/ioutil"
	"log"
	"strconv"
)

func main() {
	apimd, err := ioutil.ReadFile("./generator/API.md.tmpl")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	ioutil.WriteFile("./generator/API.md.tmpl.go", []byte("package generator\n\nconst apimdTmpl="+strconv.Quote(string(apimd))), 0777)
}
