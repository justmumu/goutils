package fileutil

import (
	"os"
	"path/filepath"
)

var (
	DefaultFilePermission   = os.FileMode(0644)
	DefaultFolderPermission = os.FileMode(0755)
)

// CleanPath cleans given path to mitigate any possible path traversal.
// it always returns an absolute path
func CleanPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Abs(path)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Abs(filepath.Join(cwd, path))
}

// CleanPathOrDefault cleans and returns the given path or returns defualt.
func CleanPathOrDefault(path, defaultPath string) string {
	if path == "" {
		return defaultPath
	}

	if val, err := CleanPath(path); err == nil {
		return val
	}

	return defaultPath
}

// SafeOpen opens a file after cleaning the path in read mode
func SafeOpen(path string) (*os.File, error) {
	abs, err := CleanPath(path)
	if err != nil {
		return nil, err
	}

	return os.Open(abs)
}

// SafeCreate creates the given path of file by cleaning path
func SafeCreate(path string) (*os.File, error) {
	abs, err := CleanPath(path)
	if err != nil {
		return nil, err
	}

	if err := CreateMissingDirs(abs); err != nil {
		return nil, err
	}

	return os.Create(abs)
}

// SafeOpenAppend opens a file after cleaning the path
// in append mode and creates any missing directories in chain /path/to/file
func SafeOpenAppend(path string) (*os.File, error) {
	abs, err := CleanPath(path)
	if err != nil {
		return nil, err
	}

	if err := CreateMissingDirs(abs); err != nil {
		return nil, err
	}

	return os.OpenFile(abs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, DefaultFilePermission)
}

// SafeOpenWrite opens a file after cleaning the path
// in write mode and creates any missing directories in chain /path/to/file
func SafeOpenWrite(path string) (*os.File, error) {
	abs, err := CleanPath(path)
	if err != nil {
		return nil, err
	}

	if err := CreateMissingDirs(abs); err != nil {
		return nil, err
	}

	return os.OpenFile(abs, os.O_CREATE|os.O_WRONLY, DefaultFilePermission)
}

// SafeWriteFile writes data to a file after cleaning the path
// in write mode and creates any missing directories in chain /path/to/file
func SafeWriteFile(path string, data []byte) error {
	abs, err := CleanPath(path)
	if err != nil {
		return err
	}

	if err := CreateMissingDirs(abs); err != nil {
		return err
	}

	return os.WriteFile(abs, data, DefaultFilePermission)
}

// SafeMkdirAll creates any missing directories in chain /path/to/file by cleaning the path
func CreateMissingDirs(path string) error {
	abs, err := CleanPath(path)
	if err != nil {
		return err
	}

	return os.MkdirAll(filepath.Dir(abs), DefaultFolderPermission)
}
