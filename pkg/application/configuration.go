package application

import (
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

// configurationLoaderInformation describes loader information structure
type configurationLoaderInformation struct {
	Extension string
	FileName  string
	Path      string
	File      string
}

// generateConfigurationLoaderInformation generates information required by the loader
func generateConfigurationLoaderInformation(application *Application) *configurationLoaderInformation {
	configFileExtension := filepath.Ext(application.configurationFile)
	configFileName := strings.TrimSuffix(application.configurationFile, configFileExtension)
	configFileExtension = configFileExtension[1:]

	viper.SetDefault("application.configuration_file_path", application.configurationPath)
	viper.SetDefault("application.configuration_file", application.configurationFile)
	viper.SetDefault("application.configuration_file_name", configFileName)

	return &configurationLoaderInformation{
		Extension: configFileExtension,
		FileName:  configFileName,
		Path:      application.configurationPath,
		File:      application.configurationFile,
	}
}
