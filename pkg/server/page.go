package server

import (
	"github.com/darki73/pac-manager/pkg/logger"
	"net/http"
	"text/template"
)

var pageTemplate = `
<!DOCTYPE html>
<html>
<head>
<title>{{.Title}}</title>
<style>
	table {
		width: 100%;
	}
	table, th, td {
		border: 1px solid;
	}
	
	td {
		text-align: center;
		vertical-align: middle;
	}
</style>
</head>
<body>
<center>
<a href="/">Home</a> | 
<a href="/proxies">Proxies</a> | 
<a href="/domains">Domains</a>
</center>
<h1>{{.Title}}</h1>
{{.Body}}
</body>
</html>
`

// getTemplate returns the page template.
func getTemplate() *template.Template {
	tpl, err := template.New("").Parse(pageTemplate)
	if err != nil {
		logger.Panicf("server:page", "failed to parse page template - %s", err.Error())
	}
	return tpl
}

// renderTemplate renders the page template.
func renderTemplate(title string, body string, writer http.ResponseWriter) {
	tpl := getTemplate()

	content := struct {
		Title string
		Body  string
	}{
		Title: title,
		Body:  body,
	}

	if err := tpl.Execute(writer, content); err != nil {
		logger.Fatalf("server:page", "failed to generate page output - %s", err.Error())
	}
}
