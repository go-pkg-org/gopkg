package archive

import (
	"encoding/json"
	"fmt"
	"github.com/go-pkg-org/gopkg/internal/pkg"
	"net/http"
)

//go:generate mockgen -destination=../archive_mock/client_mock.go -package=archive_mock . Client

// DefaultURL is the default production URL of our archive
const DefaultURL = "https://archive.gopkg.org"

// Client is an interface to dial with an archive
type Client interface {
	// GetIndex returns the up-to-date archive index
	GetIndex() (Index, error)
	// GetReleases get the available releases of given package
	GetReleases(pkgName string) (map[string][]Release, error)
	// GetLatestRelease get the latest available release of given package
	GetLatestRelease(alias, os, arch string) (pkg.File, error)
}

type client struct {
	url   string
	index Index
}

func (c *client) GetIndex() (Index, error) {
	url := fmt.Sprintf("%s/index.json", c.url)
	resp, err := http.Get(url)
	if err != nil {
		return Index{}, err
	}
	if resp.StatusCode != 200 {
		return Index{}, fmt.Errorf("error while getting index: %s", resp.Status)
	}

	var index Index
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return Index{}, fmt.Errorf("error while getting index: %s", err)
	}

	c.index = index
	return index, nil
}

func (c *client) GetReleases(pkgName string) (map[string][]Release, error) {
	// Refresh index if needed
	if len(c.index.Packages) == 0 {
		if _, err := c.GetIndex(); err != nil {
			return nil, err
		}
	}

	p, exist := c.index.Packages[pkgName]
	if !exist {
		return nil, fmt.Errorf("package %s doesn't exist", pkgName)
	}

	return p.Releases, nil
}

func (c *client) GetLatestRelease(alias, os, arch string) (pkg.File, error) {
	// Refresh index if needed
	if len(c.index.Packages) == 0 {
		if _, err := c.GetIndex(); err != nil {
			return nil, err
		}
	}

	p, exist := c.index.Packages[alias]
	if !exist {
		return nil, fmt.Errorf("package %s doesn't exist", alias)
	}

	releases := p.Releases[p.LatestRelease]

	var pkgURL string
	// Only one release in case of source package
	if len(releases) == 1 {
		pkgURL = fmt.Sprintf("%s/%s", c.url, releases[0].Path)
	} else {
		for _, release := range p.Releases[p.LatestRelease] {
			if release.OS == os && release.Arch == arch {
				pkgURL = fmt.Sprintf("%s/%s", c.url, release.Path)
				break
			}
		}
	}

	resp, err := http.Get(pkgURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error while getting last release: %s", resp.Status)
	}

	return pkg.Read(resp.Body)
}

// NewClient create a new client for an Archive
func NewClient(url string) (Client, error) {
	return &client{
		url: url,
	}, nil
}
