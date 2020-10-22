package archive

// Index represent an Archive index
// the index is used to perform packages lookup
type Index struct {
	// Packages is the list of existing packages on the archive
	Packages map[string]Package
}

// Package represent an installable package
type Package struct {
	// Releases represent the existing package releases
	Releases map[string][]Release
	// LatestRelease contains the package latest release
	LatestRelease string
	// TODO description
}

// Release represent the release of a package
type Release struct {
	OS   string
	Arch string
	Path string
}
