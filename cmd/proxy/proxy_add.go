package proxy

import (
	"github.com/darki73/pac-manager/pkg/config"
	"github.com/darki73/pac-manager/pkg/database"
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/darki73/pac-manager/pkg/pac"
	"github.com/darki73/pac-manager/pkg/prompt"
	"github.com/darki73/pac-manager/pkg/storage"
	"github.com/spf13/cobra"
	"strconv"
)

var ProxyAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a new proxy",
	Long:  "Command which allows to create new proxy server to be used in the PAC file",
	Run: func(cmd *cobra.Command, args []string) {
		configuration, err := config.Initialize()
		if err != nil {
			logger.Fatalf("proxy:add", "failed to initialize configuration: %s", err.Error())
		}

		db := database.NewDatabase(
			storage.NewStorage(configuration.GetDatabase().GetPath(), 0755),
			configuration.GetDatabase().GetName(),
		)

		if err := db.Start(); err != nil {
			logger.Fatalf("proxy:add", "failed to start database: %s", err.Error())
		}
		defer db.Stop()

		proxyType := prompt.Select("Proxy Type", "Selected proxy type is: ", []string{"HTTP", "SOCKS4", "SOCKS5"})
		proxyHost := prompt.Text("Proxy Host")
		proxyPortString := prompt.Text("Proxy Port")

		proxyPort, err := strconv.Atoi(proxyPortString)
		if err != nil {
			logger.Fatalf("proxy:add", "failed to convert proxy port to integer: %s", err.Error())
		}

		if db.ProxyCreate(proxyType, proxyHost, proxyPort) {
			logger.Infof("proxy:add", "proxy '%s %s:%d' added successfully", proxyType, proxyHost, proxyPort)

			pac.NewPac(
				storage.NewStorage(configuration.GetPac().GetPath(), 0755),
				configuration.GetPac().GetName(),
			).Generate(db.ProxyList(), db.DomainList())
		} else {
			logger.Warnf("proxy:add", "failed to add proxy '%s %s:%d'", proxyType, proxyHost, proxyPort)
		}
	},
}
