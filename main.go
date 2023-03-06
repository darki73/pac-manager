package main

import (
	"github.com/darki73/pac-manager/cmd"
	"github.com/darki73/pac-manager/cmd/domain"
	"github.com/darki73/pac-manager/cmd/proxy"
	"github.com/darki73/pac-manager/pkg/application"
	"github.com/darki73/pac-manager/pkg/logger"
)

func main() {
	app := application.NewApplication(application.Options{
		ShortDescription:       "PAC Manager",
		LongDescription:        "Proxy Auto Configuration Manager",
		HasConfiguration:       true,
		ConfigurationPath:      "/etc/pac.d",
		ConfigurationFile:      "main.yaml",
		ConfigurationEnvPrefix: "PACM",
	})

	domainCommand := domain.DomainCommand
	app.RegisterSubCommand(domainCommand, domain.DomainAddCommand)
	app.RegisterSubCommand(domainCommand, domain.DomainUpdateCommand)
	app.RegisterSubCommand(domainCommand, domain.DomainDeleteCommand)

	proxyCommand := proxy.ProxyCommand
	app.RegisterSubCommand(proxyCommand, proxy.ProxyAddCommand)
	app.RegisterSubCommand(proxyCommand, proxy.ProxyDeleteCommand)

	proxyAuthCommand := proxy.ProxyAuthCommand
	app.RegisterSubCommand(proxyAuthCommand, proxy.ProxyAuthEnableCommand)
	app.RegisterSubCommand(proxyAuthCommand, proxy.ProxyAuthDisableCommand)

	app.RegisterSubCommand(proxyCommand, proxyAuthCommand)

	app.RegisterCommand(domainCommand)
	app.RegisterCommand(proxyCommand)
	app.RegisterCommand(cmd.RunCommand)
	app.RegisterCommand(cmd.VersionCommand)

	if err := app.Start(); err != nil {
		logger.Fatalf("main", "failed to start application: %s", err.Error())
	}
}
