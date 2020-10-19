package cache

import (
	"encoding/json"
	"os"
)

// Cache represent the installed package cache
// TODO when managing archive refactor this
type Cache struct {
	Packages map[string][]string `json:"packages"`
}

// Read a cache from target path
func Read(path string) (*Cache, error) {
	var c Cache
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Cache{Packages: map[string][]string{}}, nil
		}
		return nil, err
	}

	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

// Write a cache to target path
func Write(path string, cache *Cache) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(f).Encode(*cache); err != nil {
		return err
	}
	return nil
}

// GetFiles return files associated with given package
// return nil if no such package exist
func (c *Cache) GetFiles(pkg string) []string {
	return c.Packages[pkg]
}

// AddPackage add given package alongside his files into the cache
func (c *Cache) AddPackage(pkg string, files []string) {
	c.Packages[pkg] = files
}

// RemovePackage remove given package from cache
func (c *Cache) RemovePackage(pkg string) {
	delete(c.Packages, pkg)
}
