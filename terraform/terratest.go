package tfhelpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func init() {
	// setup plugin cache to make the tests run a bit faster.
	// unless NO_TF_PLUGIN_CACHE is set
	if os.Getenv("NO_TF_PLUGIN_CACHE") == "" {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		cachedir := filepath.Join(cwd, ".plugin-cache")
		archdir := filepath.Join(cachedir, "linux_amd64")
		err = os.MkdirAll(archdir, os.FileMode(0755))
		if err != nil {
			panic(err)
		}
		os.Setenv("TF_PLUGIN_CACHE_DIR", cachedir)
	}
}

// SymlinkTerraRoot create symbolic link to terraform src from dest
func SymlinkTerraRoot(t *testing.T, srcdir, destdir string) {
	files, err := ioutil.ReadDir(srcdir)
	require.NoErrorf(t, err, "Error reading directory %s: %s", srcdir, err)

	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext == ".tf" || ext == ".zip" {
			err = os.Symlink(filepath.Join("..", srcdir, f.Name()), filepath.Join(destdir, f.Name()))
			require.NoErrorf(t, err, "Error creating sym link: %s", err)
		}
	}
}

// GetOptions create a basic Terraform Options with only directory set
func GetOptions(dir string) *terraform.Options {
	return &terraform.Options{
		TerraformDir: dir,
	}
}
