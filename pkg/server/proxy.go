package server

import (
	"encoding/base64"
	"fmt"
	"github.com/darki73/pac-manager/pkg/logger"
	"net/http"
	"strconv"
	"strings"
)

// proxyListHandler handles the proxy list page.
func (server *Server) proxyListHandler(writer http.ResponseWriter, request *http.Request) {
	proxies := server.GetDatabase().ProxyList()

	body := strings.Join([]string{
		"<a href=\"/proxies/add\">Add Proxy</a>",
		"<table>",
		"<tr>",
		"<th>Proxy</th>",
		"<th>Actions</th>",
		"</tr>",
	}, "\n")

	for _, proxy := range proxies {
		proxyHash := base64.StdEncoding.EncodeToString([]byte(proxy))
		body += fmt.Sprintf("<tr><td>%s</td><td><a href=\"/proxies/delete/%s\">Delete</a></td></tr>\n", proxy, proxyHash)
	}

	body += "</table>\n"

	renderTemplate("Proxies", body, writer)
}

// proxyAddHandler handles the proxy add page.
func (server *Server) proxyAddHandler(writer http.ResponseWriter, request *http.Request) {
	body := []string{
		"<form action=\"/proxies/save\" method=\"post\">",
		"<label for=\"proxyType\">Proxy Type:</label>",
		"<select id=\"proxyType\" name=\"proxyType\">",
		"<option value=\"HTTP\">HTTP</option>",
		"<option value=\"SOCKS4\">SOCKS4</option>",
		"<option value=\"SOCKS5\">SOCKS5</option>",
		"</select><br><br>",
		"<label for=\"proxyHost\">Proxy Host:</label>",
		"<input type=\"text\" id=\"proxyHost\" name=\"proxyHost\"><br><br>",
		"<label for=\"proxyPort\">Proxy Port:</label>",
		"<input type=\"text\" id=\"proxyPort\" name=\"proxyPort\"><br><br>",
		"<input type=\"submit\" value=\"Submit\">",
		"</form>",
	}

	renderTemplate("Add Proxy", strings.Join(body, "\n"), writer)
}

// proxyDeleteHandler handles the proxy delete page.
func (server *Server) proxyDeleteHandler(writer http.ResponseWriter, request *http.Request) {
	proxy := strings.TrimPrefix(request.URL.Path, "/proxies/delete/")
	proxyBytes, err := base64.StdEncoding.DecodeString(proxy)

	if err != nil {
		logger.Errorf("server:proxy:delete", "failed to decode proxy hash: %s", err.Error())
		http.Redirect(writer, request, "/proxies", http.StatusSeeOther)
		return
	}

	proxy = string(proxyBytes)
	proxyType, proxyHost, proxyPort := parseProxyString(proxy)

	if !server.GetDatabase().ProxyExists(proxyType, proxyHost, proxyPort) {
		logger.Errorf("server:proxy:delete", "proxy does not exist: %s (%s | %s | %d)", proxy, proxyType, proxyHost, proxyPort)
		http.Redirect(writer, request, "/proxies", http.StatusSeeOther)
		return
	}

	if server.GetDatabase().ProxyDelete(proxyType, proxyHost, proxyPort) {
		server.generatePacFile()
	}

	http.Redirect(writer, request, "/proxies", http.StatusFound)
}

// proxySaveHandler handles the proxy save page.
func (server *Server) proxySaveHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	proxyType := request.FormValue("proxyType")
	proxyHost := request.FormValue("proxyHost")
	proxyPort := request.FormValue("proxyPort")

	proxyPortInt, _ := strconv.Atoi(proxyPort)

	if !server.GetDatabase().ProxyExists(proxyType, proxyHost, proxyPortInt) {
		if server.GetDatabase().ProxyCreate(proxyType, proxyHost, proxyPortInt) {
			server.generatePacFile()
		}
	}

	http.Redirect(writer, request, "/proxies", http.StatusFound)
}

// parseProxyString parses a proxy string into its host, port and port as an int.
func parseProxyString(proxy string) (string, string, int) {
	parts := strings.Split(proxy, " ")
	proxyType := parts[0]
	parts = strings.Split(parts[1], ":")
	proxyHost := parts[0]
	proxyPort, err := strconv.Atoi(parts[1])

	if err != nil {
		logger.Errorf("server:proxy:parse", "failed to parse proxy port: %s", err.Error())
		return "", "", 0
	}

	return proxyType, proxyHost, proxyPort
}
