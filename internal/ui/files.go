package ui

import (
	"path"
	"strings"
)

type FilesPane struct {
	Items       []FileRow
	SearchQuery string
	Cursor      CursorPosition
	Selected    int
	Expanded    map[string]bool
}

type FileRow struct {
	Path  string
	Name  string
	Depth int
	IsDir bool
}

func NewFilesPane(paths []string) FilesPane {
	rows := fileRowsFromPaths(paths)
	return FilesPane{
		Items:    rows,
		Expanded: expandedDirs(rows),
	}
}

func fileRowsFromPaths(paths []string) []FileRow {
	if len(paths) == 0 {
		paths = []string{"go.mod", "cmd/", "internal/"}
	}

	rows := make([]FileRow, 0, len(paths))
	seenDirs := make(map[string]bool)
	for _, filePath := range paths {
		cleanPath := strings.TrimSuffix(path.Clean(filePath), "/")
		if cleanPath == "." {
			continue
		}

		isDir := strings.HasSuffix(filePath, "/")
		parts := strings.Split(cleanPath, "/")
		for depth := range len(parts) - 1 {
			dirPath := strings.Join(parts[:depth+1], "/")
			if seenDirs[dirPath] {
				continue
			}

			seenDirs[dirPath] = true
			rows = append(rows, FileRow{
				Path:  dirPath,
				Name:  parts[depth],
				Depth: depth,
				IsDir: true,
			})
		}

		if isDir && seenDirs[cleanPath] {
			continue
		}

		if isDir {
			seenDirs[cleanPath] = true
		}

		rows = append(rows, FileRow{
			Path:  cleanPath,
			Name:  path.Base(cleanPath),
			Depth: strings.Count(cleanPath, "/"),
			IsDir: isDir,
		})
	}

	return rows
}

func expandedDirs(rows []FileRow) map[string]bool {
	expanded := make(map[string]bool)
	for _, row := range rows {
		if row.IsDir {
			expanded[row.Path] = true
		}
	}

	return expanded
}
