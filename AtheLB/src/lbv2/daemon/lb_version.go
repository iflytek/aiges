package daemon

import (
	"log"
	"os"
	"strings"
	"text/template"
)

const (
	verFormat string = `{{printf "Ver:\t\t%v\nApiVer:\t\t%v\nDate:\t\t%v\nName:\t\t%v\nAuthor:\t\t%v\nRepo:\t\t%v\ncomments:\t%v\n" .Ver .ApiVer .Date .Name .Author .Repo .Comments}}`
)

type metaData struct {
	Ver      interface{}
	ApiVer   interface{}
	Date     interface{}
	Name     interface{}
	Author   interface{}
	Repo     interface{}
	Comments interface{}
}

func ver() {
	templateInst := template.New("ver")

	templateInstContent, templateInstContentErr := templateInst.Parse(verFormat)
	if nil != templateInstContentErr {
		log.Fatalf("templateInstContentErr:%v", templateInstContentErr)
	}

	templateInstContentExecErr := templateInstContent.Execute(
		os.Stdout,
		metaData{
			Ver:      lbVersion,
			ApiVer:   lbApiVersion,
			Date:     date,
			Name:     name,
			Author:   author,
			Repo:     repo,
			Comments: comments,
		})

	if nil != templateInstContentExecErr {
		log.Fatal(templateInstContentExecErr)
	}
}
func Version() {
	if strings.Contains(strings.ToUpper(strings.Join(os.Args, "|")), "-V") {
		std.Println("-----------------------------------------")
		ver()
		std.Println("-----------------------------------------")
		os.Exit(0)
	}
}
