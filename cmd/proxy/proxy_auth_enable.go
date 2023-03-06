package proxy

import (
	"fmt"
	"github.com/darki73/pac-manager/pkg/config"
	"github.com/darki73/pac-manager/pkg/database"
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/darki73/pac-manager/pkg/prompt"
	"github.com/darki73/pac-manager/pkg/storage"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var ProxyAuthEnableCommand = &cobra.Command{
	Use:   "enable",
	Short: "Enable the authentication for the proxy",
	Long:  "Command which allows to enable the authentication for the proxy",
	Run: func(cmd *cobra.Command, args []string) {
		configuration, err := config.Initialize()
		if err != nil {
			logger.Fatalf("proxy:auth:enable", "failed to initialize configuration: %s", err.Error())
		}

		db := database.NewDatabase(
			storage.NewStorage(configuration.GetDatabase().GetPath(), 0755),
			configuration.GetDatabase().GetName(),
		)

		if err := db.Start(); err != nil {
			logger.Fatalf("proxy:auth:enable", "failed to start database: %s", err.Error())
		}
		defer db.Stop()

		proxyHost := prompt.Text("Proxy Host")

		proxyHostsRaw := db.ProxyFindByHost(proxyHost)

		if len(proxyHostsRaw) == 0 {
			logger.Fatalf("proxy:auth:enable", "no proxy found for host '%s'", proxyHost)
		}

		proxyHosts := make([]string, len(proxyHostsRaw))
		for i, proxy := range proxyHostsRaw {
			proxyHosts[i] = fmt.Sprintf(
				"%s %s:%d",
				proxy.Type,
				proxy.Host,
				proxy.Port,
			)
		}

		enableAuth := prompt.Select("Proxy to enable authentication", "Selected proxy to enable authentication is: ", proxyHosts)

		parts := strings.Split(enableAuth, " ")
		proxyType := parts[0]
		parts = strings.Split(parts[1], ":")
		proxyHost = parts[0]

		proxyPort, err := strconv.Atoi(parts[1])
		if err != nil {
			logger.Fatalf("proxy:auth:enable", "failed to convert port to integer: %s", err.Error())
		}

		username := prompt.Text("Username")
		password := prompt.Password("Password")

		if db.ProxyEnableAuthentication(proxyType, proxyHost, proxyPort, username, password) {
			logger.Infof("proxy:auth:enable", "authentication enabled for proxy '%s %s:%d'", proxyType, proxyHost, proxyPort)
		} else {
			logger.Fatalf("proxy:auth:enable", "failed to enable authentication for proxy '%s %s:%d'", proxyType, proxyHost, proxyPort)
		}
	},
}
