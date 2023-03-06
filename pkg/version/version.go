package version

import "time"

var (
	// version holds the current version of binary.
	version = "dev"

	// commit holds the current commit of binary.
	commit = "none"

	// buildDate holds the build date of binary.
	buildDate = "0000-00-00T00-00-00Z"

	// builtBy holds the builder name of binary.
	builtBy = "freedomcore"

	// startDate holds the start date of binary.
	startDateTime = time.Now()
)

// GetVersion returns current version of binary.
func GetVersion() string {
	return version
}

// GetCommit returns commit with which binary was built.
func GetCommit() string {
	return commit
}

// GetBuildDate returns date on which binary was built.
func GetBuildDate() string {
	return buildDate
}

// GetBuiltBy returns binary builder name.
func GetBuiltBy() string {
	return builtBy
}

// GetStartDateTime returns binary start date.
func GetStartDateTime() time.Time {
	return startDateTime
}
