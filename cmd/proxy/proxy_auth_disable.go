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

var ProxyAuthDisableCommand = &cobra.Command{
	Use:   "disable",
	Short: "Disable the authentication",
	Long:  "Command which allows to disable the authentication for the proxy server",
	Run: func(cmd *cobra.Command, args []string) {
		configuration, err := config.Initialize()
		if err != nil {
			logger.Fatalf("proxy:auth:disable", "failed to initialize configuration: %s", err.Error())
		}

		db := database.NewDatabase(
			storage.NewStorage(configuration.GetDatabase().GetPath(), 0755),
			configuration.GetDatabase().GetName(),
		)

		if err := db.Start(); err != nil {
			logger.Fatalf("proxy:auth:disable", "failed to start database: %s", err.Error())
		}
		defer db.Stop()

		proxyHost := prompt.Text("Proxy Host")

		proxyHostsRaw := db.ProxyFindByHost(proxyHost)

		if len(proxyHostsRaw) == 0 {
			logger.Fatalf("proxy:auth:disable", "no proxy found for host '%s'", proxyHost)
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

		disableAuth := prompt.Select("Proxy to disable authentication", "Selected proxy to disable authentication is: ", proxyHosts)

		parts := strings.Split(disableAuth, " ")
		proxyType := parts[0]
		parts = strings.Split(parts[1], ":")
		proxyHost = parts[0]

		proxyPort, err := strconv.Atoi(parts[1])
		if err != nil {
			logger.Fatalf("proxy:auth:disable", "failed to convert port to integer: %s", err.Error())
		}

		ensure := prompt.Confirm("Are you sure you want to disable the authentication for the proxy server?")
		if ensure {
			if db.ProxyDisableAuthentication(proxyType, proxyHost, proxyPort) {
				logger.Infof("proxy:auth:disable", "authentication disabled for proxy server '%s'", disableAuth)
			} else {
				logger.Fatalf("proxy:auth:disable", "failed to disable authentication for proxy server '%s'", disableAuth)
			}
		}
	},
}
