package domain

import (
	"github.com/darki73/pac-manager/pkg/config"
	"github.com/darki73/pac-manager/pkg/database"
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/darki73/pac-manager/pkg/pac"
	"github.com/darki73/pac-manager/pkg/prompt"
	"github.com/darki73/pac-manager/pkg/storage"
	"github.com/spf13/cobra"
)

var DomainUpdateCommand = &cobra.Command{
	Use:   "update",
	Short: "Update domain",
	Long:  "Updates domain",
	Run: func(cmd *cobra.Command, args []string) {
		configuration, err := config.Initialize()
		if err != nil {
			logger.Fatalf("domain:update", "failed to initialize configuration: %s", err.Error())
		}

		db := database.NewDatabase(
			storage.NewStorage(configuration.GetDatabase().GetPath(), 0755),
			configuration.GetDatabase().GetName(),
		)

		if err := db.Start(); err != nil {
			logger.Fatalf("domain:update", "failed to start database: %s", err.Error())
		}
		defer db.Stop()

		domainName := prompt.Text("Domain Name")

		if db.DomainExists(domainName) {
			proxify := prompt.Confirm("Route through proxy")

			if db.DomainUpdate(domainName, proxify) {
				logger.Infof("domain:update", "domain '%s' updated successfully", domainName)

				pac.NewPac(
					storage.NewStorage(configuration.GetPac().GetPath(), 0755),
					configuration.GetPac().GetName(),
				).Generate(db.ProxyList(), db.DomainList())
			} else {
				logger.Warnf("domain:update", "failed to update domain '%s'", domainName)
			}
		} else {
			logger.Warnf("domain:update", "domain '%s' does not exist", domainName)
		}
	},
}
