package proxy

import (
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/spf13/cobra"
)

var ProxyCommand = &cobra.Command{
	Use:   "proxy",
	Short: "Proxy commands",
	Long:  "Displays commands which are related to proxy servers management",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			logger.Fatalf("cmd:proxy", "error while executing help command: %s", err.Error())
		}
	},
}
