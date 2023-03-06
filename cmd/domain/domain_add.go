package domain

import (
	"github.com/darki73/pac-manager/pkg/config"
	"github.com/darki73/pac-manager/pkg/database"
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/darki73/pac-manager/pkg/pac"
	"github.com/darki73/pac-manager/pkg/prompt"
	"github.com/darki73/pac-manager/pkg/storage"
	_ "github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var DomainAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add domain",
	Long:  "Adds a new domain to the database",
	Run: func(cmd *cobra.Command, args []string) {
		configuration, err := config.Initialize()
		if err != nil {
			logger.Fatalf("domain:add", "failed to initialize configuration: %s", err.Error())
		}

		db := database.NewDatabase(
			storage.NewStorage(configuration.GetDatabase().GetPath(), 0755),
			configuration.GetDatabase().GetName(),
		)

		if err := db.Start(); err != nil {
			logger.Fatalf("domain:add", "failed to start database: %s", err.Error())
		}
		defer db.Stop()

		domainName := prompt.Text("Domain Name")
		proxify := prompt.Confirm("Route through proxy")

		if db.DomainCreate(domainName, proxify) {
			logger.Infof("domain:add", "domain '%s' added successfully", domainName)

			pac.NewPac(
				storage.NewStorage(configuration.GetPac().GetPath(), 0755),
				configuration.GetPac().GetName(),
			).Generate(db.ProxyList(), db.DomainList())
		} else {
			logger.Warnf("domain:add", "failed to add domain '%s'", domainName)
		}
	},
}
