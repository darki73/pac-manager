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

var DomainDeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete domain",
	Long:  "Deletes domain from the database",
	Run: func(cmd *cobra.Command, args []string) {
		configuration, err := config.Initialize()
		if err != nil {
			logger.Fatalf("domain:delete", "failed to initialize configuration: %s", err.Error())
		}

		db := database.NewDatabase(
			storage.NewStorage(configuration.GetDatabase().GetPath(), 0755),
			configuration.GetDatabase().GetName(),
		)

		if err := db.Start(); err != nil {
			logger.Fatalf("domain:delete", "failed to start database: %s", err.Error())
		}
		defer db.Stop()

		domainName := prompt.Text("Domain Name")

		if db.DomainExists(domainName) {
			delete := prompt.Confirm("Delete domain")
			if delete {
				if db.DomainDelete(domainName) {
					logger.Infof("domain:delete", "domain '%s' deleted successfully", domainName)

					pac.NewPac(
						storage.NewStorage(configuration.GetPac().GetPath(), 0755),
						configuration.GetPac().GetName(),
					).Generate(db.ProxyList(), db.DomainList())
				} else {
					logger.Warnf("domain:delete", "failed to delete domain '%s'", domainName)
				}
			}
		} else {
			logger.Warnf("domain:delete", "domain '%s' does not exist", domainName)
		}
	},
}
