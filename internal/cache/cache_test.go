package cache

import (
	"github.com/go-pkg-org/gopkg/internal/archive"
	"github.com/go-pkg-org/gopkg/internal/archive_mock"
	"github.com/go-pkg-org/gopkg/internal/config"
	"github.com/go-pkg-org/gopkg/internal/pkg"
	"github.com/go-pkg-org/gopkg/internal/pkg_mock"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"runtime"
	"testing"
)

func TestCache_InstallPkg_AlreadyInstalled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := pkg_mock.NewMockFile(ctrl)
	p.EXPECT().Metadata().Return(pkg.Meta{Alias: "foo/bar"}, nil)

	cache := cache{
		Packages: map[string][]string{"foo/bar": {}},
	}

	if _, err := cache.installPkg(p); err != ErrPackageAlreadyInstalled {
		t.FailNow()
	}
}

func TestCache_InstallPkg_WrongTarget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := pkg_mock.NewMockFile(ctrl)
	p.EXPECT().Metadata().Return(pkg.Meta{
		Alias:      "foo/bar",
		TargetOS:   "foo",
		TargetArch: "bar",
		Main:       "main.go",
		BinName:    "foo-bar",
	}, nil)

	cache := cache{}

	if _, err := cache.installPkg(p); err != ErrWrongTarget {
		t.Errorf("installPkg should have failed with ErrWrongTarget")
	}
}

func TestCache_InstallPkg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	binDir, _ := ioutil.TempDir("", "")
	f, _ := ioutil.TempFile("", "")

	p := pkg_mock.NewMockFile(ctrl)
	p.EXPECT().Metadata().Return(pkg.Meta{
		Alias:      "foo/bar",
		TargetOS:   runtime.GOOS,
		TargetArch: runtime.GOARCH,
		Main:       "main.go",
		BinName:    "foo-bar",
	}, nil)
	p.EXPECT().Files().Return(map[string][]byte{"hello": []byte("world")})

	cache := cache{
		Packages:  map[string][]string{},
		cacheFile: f.Name(),
		conf: &config.Config{
			BinDir: binDir,
		},
	}

	if _, err := cache.installPkg(p); err != nil {
		t.Errorf("installPkg has failed: %s", err)
	}

	if len(cache.Packages) != 1 {
		t.Errorf("wrong number of packages: %d", len(cache.Packages))
	}
}

func TestCache_ListPackages_Archive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	arc := archive_mock.NewMockClient(ctrl)
	arc.EXPECT().GetIndex().Return(archive.Index{
		Packages: map[string]archive.Package{"foo/bar": {}},
	}, nil)

	cache := cache{arcClient: arc}

	pkgs, err := cache.ListPackages(false)
	if err != nil {
		t.Error(err)
	}

	if len(pkgs) != 1 {
		t.Errorf("wrong number of packages")
	}
}

func TestCache_ListPackages_Local(t *testing.T) {
	cache := &cache{
		Packages: map[string][]string{"foo/bar": {}, "local/host": {}},
	}

	pkgs, err := cache.ListPackages(true)
	if err != nil {
		t.Error(err)
	}

	if len(pkgs) != 2 {
		t.Errorf("wrong number of packages")
	}
}

func TestCache_RemovePkg_NotInstalled(t *testing.T) {
	cache := cache{
		Packages: map[string][]string{},
	}

	if err := cache.RemovePkg("foo/bar"); err == nil {
		t.Error("should have failed")
	}
}

func TestCache_RemovePkg(t *testing.T) {
	f, _ := ioutil.TempFile("", "")

	cache := cache{
		Packages:  map[string][]string{"foo/bar": {}},
		cacheFile: f.Name(),
	}

	if err := cache.RemovePkg("foo/bar"); err != nil {
		t.Error(err)
	}
}
