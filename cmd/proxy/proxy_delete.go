package proxy

import (
	"fmt"
	"github.com/darki73/pac-manager/pkg/config"
	"github.com/darki73/pac-manager/pkg/database"
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/darki73/pac-manager/pkg/pac"
	"github.com/darki73/pac-manager/pkg/prompt"
	"github.com/darki73/pac-manager/pkg/storage"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var ProxyDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a proxy",
	Long:  "Command which allows to delete a proxy server",
	Run: func(cmd *cobra.Command, args []string) {
		configuration, err := config.Initialize()
		if err != nil {
			logger.Fatalf("proxy:delete", "failed to initialize configuration: %s", err.Error())
		}

		db := database.NewDatabase(
			storage.NewStorage(configuration.GetDatabase().GetPath(), 0755),
			configuration.GetDatabase().GetName(),
		)

		if err := db.Start(); err != nil {
			logger.Fatalf("proxy:delete", "failed to start database: %s", err.Error())
		}
		defer db.Stop()

		proxyHost := prompt.Text("Proxy Host")

		proxyHostsRaw := db.ProxyFindByHost(proxyHost)

		if len(proxyHostsRaw) == 0 {
			logger.Fatalf("proxy:delete", "no proxy found for host '%s'", proxyHost)
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

		deleteProxy := prompt.Select("Proxy to delete", "Selected proxy to delete is: ", proxyHosts)

		parts := strings.Split(deleteProxy, " ")
		proxyType := parts[0]
		parts = strings.Split(parts[1], ":")
		proxyHost = parts[0]
		proxyPort, err := strconv.Atoi(parts[1])

		if err != nil {
			logger.Fatalf("proxy:delete", "failed to convert proxy port to integer: %s", err.Error())
		}

		if db.ProxyDelete(proxyType, proxyHost, proxyPort) {
			logger.Infof("proxy:delete", "proxy '%s %s:%d' deleted successfully", proxyType, proxyHost, proxyPort)

			pac.NewPac(
				storage.NewStorage(configuration.GetPac().GetPath(), 0755),
				configuration.GetPac().GetName(),
			).Generate(db.ProxyList(), db.DomainList())
		} else {
			logger.Warnf("proxy:delete", "failed to delete proxy '%s %s:%d'", proxyType, proxyHost, proxyPort)
		}
	},
}
