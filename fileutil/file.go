package fileutil

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"os"
)

// FileExists checks if the given path is a file and if it is exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		return false
	}

	return !info.IsDir()
}

// FolderExists checks if the given path is a folder and if it is exists
func FolderExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		return false
	}

	return info.IsDir()
}

// FileOrFolderExists checks if the file/folder exists
func FileOrFolderExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// DownloadFile to specified path
func DownloadFile(path string, url string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := SafeCreate(path)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}

// CreateFolders in the list
func CreateFolders(paths ...string) error {
	for _, path := range paths {
		if err := CreateFolder(path); err != nil {
			return err
		}
	}
	return nil
}

// CreateFolder path
func CreateFolder(path string) error {
	abs, err := CleanPath(path)
	if err != nil {
		return err
	}
	return os.MkdirAll(abs, DefaultFolderPermission)
}

// GetTempFileName generate a temporary file name
func GetTempFileName() (string, error) {
	tmpfile, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}

	tmpFileName := tmpfile.Name()

	if err := tmpfile.Close(); err != nil {
		return tmpFileName, err
	}

	if err := os.RemoveAll(tmpFileName); err != nil {
		return tmpFileName, err
	}
	return tmpFileName, nil
}

// CopyFile from source to destination
func CopyFile(src, dst string) error {
	srcAbs, err := CleanPath(src)
	if err != nil {
		return err
	}

	if !FileExists(srcAbs) {
		return errors.New("source file doesn't exist")
	}

	srcFile, err := SafeOpen(srcAbs)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := SafeCreate(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Sync()
}
