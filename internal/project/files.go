package project

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func ListFiles(root string) ([]string, error) {
	files := make([]string, 0)

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := entry.Name()
		if entry.IsDir() {
			switch name {
			case ".git", "vendor", "node_modules":
				return filepath.SkipDir
			}

			return nil
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		files = append(files, filepath.ToSlash(relativePath))
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

func ReadFile(root, path string) (string, error) {
	content, err := os.ReadFile(filepath.Join(root, path))

	return string(content), err
}
