package ui

type EditorPane struct {
	Content        string
	Cursor         CursorPosition
	Path           string
	IsSearchActive bool
	SearchQuery    string
	searchMatchEnd int
	IsModified     bool
}

func NewEditorPane() EditorPane {
	return EditorPane{
		Content: "",
		Cursor: CursorPosition{
			Row:    0,
			Column: 0,
		},
		Path:           "",
		IsSearchActive: false,
		SearchQuery:    "",
		searchMatchEnd: -1,
		IsModified:     false,
	}
}
