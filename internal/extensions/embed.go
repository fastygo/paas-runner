package extensions

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed *.yml
var embedded embed.FS

func List() ([]string, error) {
	entries, err := fs.ReadDir(embedded, ".")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if len(name) > 4 && name[len(name)-4:] == ".yml" {
			names = append(names, name[:len(name)-4])
		}
	}

	return names, nil
}

func Read(name string) ([]byte, error) {
	data, err := embedded.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("read embedded extension %q: %w", name, err)
	}

	return data, nil
}
