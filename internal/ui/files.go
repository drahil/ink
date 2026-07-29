package ui

type FilesPane struct {
	Items       []string
	SearchQuery string
	Cursor      CursorPosition
	Selected    int
}

func NewFilesPane(items []string) FilesPane {
	if len(items) == 0 {
		items = []string{"go.mod", "cmd/", "internal/"}
	}

	return FilesPane{
		Items: items,
	}
}
