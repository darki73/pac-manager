package cmd

import (
	"github.com/darki73/pac-manager/pkg/config"
	"github.com/darki73/pac-manager/pkg/database"
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/darki73/pac-manager/pkg/server"
	"github.com/darki73/pac-manager/pkg/storage"
	"github.com/spf13/cobra"
)

var RunCommand = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	Long:  "Run the application",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		configuration, err := config.Initialize()
		if err != nil {
			logger.Fatalf("cmd:run", "failed to initialize configuration: %s", err.Error())
		}

		db := database.NewDatabase(
			storage.NewStorage(configuration.GetDatabase().GetPath(), 0755),
			configuration.GetDatabase().GetName(),
		)

		if err := db.Start(); err != nil {
			logger.Fatalf("cmd:run", "failed to start database: %s", err.Error())
		}
		defer db.Stop()

		server.NewServer(
			configuration,
			db,
		).Run()
	},
}
