package cache

import (
	"encoding/json"
	"os"
)

type Cache struct {
	Packages map[string][]string `json:"packages"`
}

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

func (c *Cache) GetFiles(pkg string) []string {
	return c.Packages[pkg]
}

func (c *Cache) AddPackage(pkg string, files []string) {
	c.Packages[pkg] = files
}
