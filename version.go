package healthfx

import "fmt"

// Version is the current version of the application.
type Version struct {
	Version   string // semver
	ReleaseID int    // build number
	BuildDate string // ISO 8601
	GitCommit string // commit hash
	GoVersion string // go version
}

func (v Version) String() string {
	return fmt.Sprintf("%s (build %d, commit %s)", v.Version, v.ReleaseID, v.GitCommit)
}
