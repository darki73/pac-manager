package application

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Application represents structure of application.
type Application struct {
	callSign               string
	name                   string
	shortDescription       string
	longDescription        string
	hasConfiguration       bool
	configurationPath      string
	configurationFile      string
	configurationEnvPrefix string
}

// Options represents structure of application options.
type Options struct {
	ShortDescription       string
	LongDescription        string
	HasConfiguration       bool
	ConfigurationPath      string
	ConfigurationFile      string
	ConfigurationEnvPrefix string
}

// rootCommand represents the base command when called without any subcommands.
var rootCommand *cobra.Command

// NewApplication creates new instance of application
func NewApplication(options Options) *Application {
	application := &Application{
		shortDescription:       options.ShortDescription,
		longDescription:        options.LongDescription,
		hasConfiguration:       options.HasConfiguration,
		configurationPath:      options.ConfigurationPath,
		configurationFile:      options.ConfigurationFile,
		configurationEnvPrefix: options.ConfigurationEnvPrefix,
	}
	executableName := application.executableName()
	application.callSign = executableName
	application.name = executableName

	application.initializeRootCommand()

	if application.hasConfiguration {
		application.initializeRootFlags()
	}

	return application
}

// Start starts application execution
func (application *Application) Start() error {
	if application.hasConfiguration {
		cobra.OnInitialize(application.bootstrapConfigurationManager)
	}
	return rootCommand.Execute()
}

// FlagRegistrar exposes interface to work with application flags
func (application *Application) FlagRegistrar() *pflag.FlagSet {
	return rootCommand.PersistentFlags()
}

// RegisterCommand registers new application command
func (application *Application) RegisterCommand(command *cobra.Command) {
	rootCommand.AddCommand(command)
}

// RegisterSubCommand registers new application sub-command
func (application *Application) RegisterSubCommand(parentCommand *cobra.Command, command *cobra.Command) {
	parentCommand.AddCommand(command)
}

// initializeRootCommand initializes root command for application
func (application *Application) initializeRootCommand() {
	rootCommand = &cobra.Command{
		Use:               application.callSign,
		Short:             application.shortDescription,
		Long:              application.longDescription,
		DisableAutoGenTag: true,
	}
}

// initializeRootFlags initializes root flags for application
func (application *Application) initializeRootFlags() {
	application.FlagRegistrar().StringVarP(
		&application.configurationPath,
		"config-path",
		"p",
		application.configurationPath,
		"Full path to configuration directory",
	)

	application.FlagRegistrar().StringVarP(
		&application.configurationFile,
		"config",
		"c",
		application.configurationFile,
		"Name of configuration file to be used",
	)
}

// bootstrapConfigurationManager bootstraps configuration manager and registers Viper handler
func (application *Application) bootstrapConfigurationManager() {
	if application.configurationEnvPrefix != "" {
		viper.SetEnvPrefix(application.configurationEnvPrefix)
	}
	viper.AutomaticEnv()

	validPaths := application.updateConfigurationPaths()
	if !validPaths {
		logger.Trace("application", "application is expected to have configuration, but provided information is invalid")
		return
	} else {
		loaderInformation := generateConfigurationLoaderInformation(application)
		viper.SetConfigName(loaderInformation.FileName)
		viper.SetConfigType(loaderInformation.Extension)
		viper.AddConfigPath(loaderInformation.Path)

		err := viper.ReadInConfig()
		if err != nil {
			switch e := err.(type) {
			case viper.ConfigFileNotFoundError:
				logger.Errorf("application", "unable to find configuration file - %s/%s.", loaderInformation.Path, loaderInformation.File)
			case viper.ConfigParseError:
				logger.Fatalf("application", "failed to parse configuration file - %s", e.Error())
			default:
				logger.Fatalf("application", "unknown error - %s", e.Error())
			}
		}
	}
}

// updateConfigurationPaths updates configuration paths (trimming spaces, etc)
func (application *Application) updateConfigurationPaths() bool {
	confPath := strings.TrimSpace(application.configurationPath)
	if confPath == "" {
		return false
	} else {
		application.configurationPath = confPath
	}

	confFile := strings.TrimSpace(application.configurationFile)
	if confFile == "" {
		return false
	} else {
		application.configurationFile = confFile
	}

	return true
}

// executableName returns name of executable
func (application *Application) executableName() string {
	executablePath := os.Args[0]
	return filepath.Base(executablePath)
}
