package ui

type EditorPane struct {
	Content string
	Cursor  CursorPosition
	Path    string
}

func NewEditorPane() EditorPane {
	return EditorPane{
		Content: "<?php\n\necho 'hello';\n",
		Cursor: CursorPosition{
			Row:    0,
			Column: 0,
		},
		Path: "",
	}
}
