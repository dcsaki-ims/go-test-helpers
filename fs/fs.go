package fshelpers

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// CreateTempDir creates a temporary directory in the parent 'dir'
// or if 'dir' is empty the system default tmp directory.
// The created directory will use the specified prefix or the test name as a prefix
func CreateTempDir(t *testing.T, dir, prefix string) string {
	t.Helper()
	if prefix == "" {
		prefix = t.Name() + "-"
	}
	tempdir, err := ioutil.TempDir(dir, prefix)
	require.NoErrorf(t, err, "Temp directory creation failed: %v", err)
	return tempdir
}

// DeleteTempDir delete the directory
func DeleteTempDir(t *testing.T, tempdir string) {
	t.Helper()
	err := os.RemoveAll(tempdir)
	require.NoErrorf(t, err, "Temp directory deletion failed: %v", err)
}
