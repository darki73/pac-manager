package server

import (
	"fmt"
	"net/http"
	"strings"
)

// domainListHandler renders the page for listing domains.
func (server *Server) domainListHandler(writer http.ResponseWriter, request *http.Request) {
	domains := server.GetDatabase().DomainList()

	body := strings.Join([]string{
		"<a href=\"/domains/add\">Add Domain</a>",
		"<table>",
		"<tr>",
		"<th>Domain</th>",
		"<th>Should Proxy</th>",
		"<th>Actions</th>",
		"</tr>",
	}, "\n")

	for domain, shouldProxy := range domains {
		body += fmt.Sprintf("<tr><td>%s</td><td>%t</td><td><a href=\"/domains/edit/%s\">Edit</a> | <a href=\"/domains/delete/%s\">Delete</a></td></tr>\n", domain, shouldProxy, domain, domain)
	}

	body += "</table>\n"

	renderTemplate("Domains", body, writer)
}

// domainEditHandler renders the page for editing a domain.
func (server *Server) domainEditHandler(writer http.ResponseWriter, request *http.Request) {
	domain := strings.TrimPrefix(request.URL.Path, "/domains/edit/")

	if !server.GetDatabase().DomainExists(domain) {
		http.Redirect(writer, request, "/domains", http.StatusNotFound)
		return
	}

	_, name, shouldProxy := server.GetDatabase().DomainGet(domain)

	body := []string{
		"<form action=\"/domains/save\" method=\"post\">",
		"<label for=\"domain\">Domain:</label>",
		"<input type=\"text\" id=\"domain\" name=\"domain\" value=\"" + name + "\"><br><br>",
		"<label for=\"shouldProxy\">Should Proxy:</label>",
	}

	if shouldProxy {
		body = append(
			body,
			"<input type=\"checkbox\" id=\"shouldProxy\" name=\"shouldProxy\" value=\"true\" checked=\"true\"><br/><br/>",
		)
	} else {
		body = append(
			body,
			"<input type=\"checkbox\" id=\"shouldProxy\" name=\"shouldProxy\" value=\"true\"><br/><br/>",
		)
	}

	body = append(body, "<input type=\"submit\" value=\"Submit\">")
	body = append(body, "</form>")
	renderTemplate("Edit Domain", strings.Join(body, "\n"), writer)
}

// domainDeleteHandler deletes a domain.
func (server *Server) domainDeleteHandler(writer http.ResponseWriter, request *http.Request) {
	domain := strings.TrimPrefix(request.URL.Path, "/domains/delete/")

	if !server.GetDatabase().DomainExists(domain) {
		http.Redirect(writer, request, "/domains", http.StatusNotFound)
		return
	}

	if server.GetDatabase().DomainDelete(domain) {
		server.generatePacFile()
	}

	http.Redirect(writer, request, "/domains", http.StatusSeeOther)
}

// domainAddHandler renders the page for adding a domain.
func (server *Server) domainAddHandler(writer http.ResponseWriter, request *http.Request) {
	body := strings.Join([]string{
		"<form action=\"/domains/save\" method=\"post\">",
		"<label for=\"domain\">Domain:</label>",
		"<input type=\"text\" id=\"domain\" name=\"domain\"><br><br>",
		"<label for=\"shouldProxy\">Should Proxy:</label>",
		"<input type=\"checkbox\" id=\"shouldProxy\" name=\"shouldProxy\" value=\"true\"><br><br>",
		"<input type=\"submit\" value=\"Submit\">",
		"</form>",
	}, "\n")

	renderTemplate("Create Domain", body, writer)
}

// domainSaveHandler saves/updates a domain.
func (server *Server) domainSaveHandler(writer http.ResponseWriter, request *http.Request) {
	domain := request.FormValue("domain")
	shouldProxyRaw := request.FormValue("shouldProxy")

	shouldProxy := false
	if shouldProxyRaw == "true" {
		shouldProxy = true
	}

	if server.GetDatabase().DomainExists(domain) {
		server.GetDatabase().DomainUpdate(domain, shouldProxy)
	} else {
		server.GetDatabase().DomainCreate(domain, shouldProxy)
	}

	server.generatePacFile()

	http.Redirect(writer, request, "/domains", http.StatusSeeOther)
}
