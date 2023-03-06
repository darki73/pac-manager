package cmd

import (
	"html/template"
	"os"
	"runtime"

	"github.com/darki73/pac-manager/pkg/logger"
	"github.com/darki73/pac-manager/pkg/version"
	"github.com/spf13/cobra"
)

var versionTemplate = `Version:      {{ .Version }}
SHA Commit:   {{ .Commit }}
Go version:   {{ .GoVersion }}
Built On:     {{ .BuildDate }}
Built By:     {{ .BuiltBy }}
OS/Arch:      {{ .Os }}/{{ .Arch }}
`

var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "Display application version",
	Long:  "Get version of installed application binary",
	Run: func(cmd *cobra.Command, args []string) {
		tpl, err := template.New("").Parse(versionTemplate)
		if err != nil {
			logger.Panicf("application:version", "failed to parse version template - %s", err.Error())
		}
		v := struct {
			Version   string
			Commit    string
			GoVersion string
			BuildDate string
			BuiltBy   string
			Os        string
			Arch      string
		}{
			Version:   version.GetVersion(),
			Commit:    version.GetCommit(),
			GoVersion: runtime.Version(),
			BuildDate: version.GetBuildDate(),
			BuiltBy:   version.GetBuiltBy(),
			Os:        runtime.GOOS,
			Arch:      runtime.GOARCH,
		}
		err = tpl.Execute(os.Stdout, v)
		if err != nil {
			logger.Fatalf("application:version", "failed to generate version output - %s", err.Error())
		}
	},
}
