package proxy

import (
	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/spf13/cobra"
)

var ProxyAuthCommand = &cobra.Command{
	Use:   "auth",
	Short: "Proxy Authentication Management",
	Long:  "Command which allows to manage proxy authentication settings",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			logger.Fatalf("cmd:proxy", "error while executing help command: %s", err.Error())
		}
	},
}
