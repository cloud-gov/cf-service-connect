// Package version holds the plugin's version in one place, so that the version
// reported to `cf plugins` and the version sent to the CF API in the User-Agent
// header cannot drift apart.
package version

import "fmt"

// The plugin version, as reported by `cf plugins`.
const (
	Major = 2
	Minor = 0
	Build = 1
)

// String renders the version as "major.minor.build".
func String() string {
	return fmt.Sprintf("%d.%d.%d", Major, Minor, Build)
}
