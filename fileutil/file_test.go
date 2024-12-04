package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileOrFolderExists(t *testing.T) {
	tests := map[string]bool{
		"file.go":     true,
		"aaa.bbb":     false,
		".":           true,
		"../fileutil": true,
		"aabb":        false,
	}
	for fpath, mustExist := range tests {
		exist := FileOrFolderExists(fpath)
		require.Equalf(t, mustExist, exist, "invalid \"%s\": %v", fpath, exist)
	}
}

func TestFileExists(t *testing.T) {
	tests := map[string]bool{
		"file.go": true,
		"aaa.bbb": false,
		"/":       false,
	}
	for fpath, mustExist := range tests {
		exist := FileExists(fpath)
		require.Equalf(t, mustExist, exist, "invalid \"%s\": %v", fpath, exist)
	}
}

func TestFolderExists(t *testing.T) {
	tests := map[string]bool{
		".":           true,
		"../fileutil": true,
		"aabb":        false,
	}
	for fpath, mustExist := range tests {
		exist := FolderExists(fpath)
		require.Equalf(t, mustExist, exist, "invalid \"%s\"", fpath)
	}
}

func TestDownloadFile(t *testing.T) {
	// attempt to download http://ipv4.download.thinkbroadband.com/5MB.zip to temp folder
	tmpfile, err := os.CreateTemp("", "")
	require.Nil(t, err, "couldn't create folder: %s", err)
	fname := tmpfile.Name()

	os.Remove(fname)

	err = DownloadFile(fname, "http://ipv4.download.thinkbroadband.com/5MB.zip")
	require.Nil(t, err, "couldn't download file: %s", err)

	require.True(t, FileExists(fname), "file \"%s\" doesn't exists", fname)

	// remove the downloaded file
	os.Remove(fname)
}

func tmpFolderName(s string) string {
	return filepath.Join(os.TempDir(), s)
}

func TestCreateFolders(t *testing.T) {
	tests := []string{
		tmpFolderName("a"),
		tmpFolderName("b"),
	}
	err := CreateFolders(tests...)
	require.Nil(t, err, "couldn't download file: %s", err)

	for _, folder := range tests {
		fexists := FolderExists(folder)
		require.True(t, fexists, "folder %s doesn't exist", folder)
	}

	// remove folders
	for _, folder := range tests {
		os.Remove(folder)
	}
}

func TestCreateFolder(t *testing.T) {
	tst := tmpFolderName("a")
	err := CreateFolder(tst)
	require.Nil(t, err, "couldn't download file: %s", err)

	fexists := FolderExists(tst)
	require.True(t, fexists, "folder %s doesn't exist", fexists)

	os.Remove(tst)
}

func TestCopyFile(t *testing.T) {
	fileContent := `test
	test1
	test2`
	f, err := os.CreateTemp("", "")
	require.Nil(t, err, "couldn't create file: %s", err)
	fname := f.Name()
	_, _ = f.Write([]byte(fileContent))
	f.Close()
	defer os.Remove(fname)
	fnameCopy := fmt.Sprintf("%s-copy", f.Name())
	err = CopyFile(fname, fnameCopy)
	require.Nil(t, err, "couldn't copy file: %s", err)
	require.True(t, FileExists(fnameCopy), "file \"%s\" doesn't exists", fnameCopy)
	os.Remove(fnameCopy)
}

func TestGetTempFileName(t *testing.T) {
	fname, _ := GetTempFileName()
	defer os.Remove(fname)
	require.NotEmpty(t, fname)
}
