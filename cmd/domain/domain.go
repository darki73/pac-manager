package domain

import (
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/spf13/cobra"
)

var DomainCommand = &cobra.Command{
	Use:   "domain",
	Short: "Domain commands",
	Long:  "Displays commands which are related to domain management",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			logger.Fatalf("cmd:domain", "error while executing help command: %s", err.Error())
		}
	},
}
