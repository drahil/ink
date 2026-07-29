package project

type Store struct {
	Root string
}

func NewStore(root string) Store {
	return Store{Root: root}
}

func (s Store) ListFiles() ([]string, error) {
	return ListFiles(s.Root)
}

func (s Store) ReadFile(path string) (string, error) {
	return ReadFile(s.Root, path)
}

func (s Store) WriteFile(path, content string) error {
	return WriteFile(s.Root, path, content)
}
