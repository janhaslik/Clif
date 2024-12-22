package filesystem

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type FileBrowser struct {
	currentPath      string
	currentPathEntry os.DirEntry
	dirEntries       []os.DirEntry
}

type dirEntry struct {
	name  string
	isDir bool
}

func (d *dirEntry) Name() string {
	return d.name
}

func (d *dirEntry) IsDir() bool {
	return d.isDir
}

func (d *dirEntry) Type() os.FileMode {
	if d.isDir {
		return os.ModeDir
	}
	return 0
}

func (d *dirEntry) Info() (os.FileInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func NewFileBrowser(initialDir string) *FileBrowser {
	if _, err := os.Stat(initialDir); err != nil {
		log.Fatalf("Error: Initial directory %s does not exist or cannot be accessed: %v", initialDir, err)
	}

	initialDirEntry := &dirEntry{
		name:  filepath.Base(initialDir),
		isDir: true,
	}

	f := &FileBrowser{
		currentPath:      initialDir,
		currentPathEntry: initialDirEntry,
	}
	f.dirEntries = f.listDir()
	return f
}

func (f *FileBrowser) CurrentPath() string {
	return f.currentPath
}

func (f *FileBrowser) CurrentPathEntry() os.DirEntry {
	return f.currentPathEntry
}

func (f *FileBrowser) DirEntries() []os.DirEntry {
	return f.dirEntries
}

func (f *FileBrowser) SetCurrentDirEntry(name string) error {
	entries := f.listDir()

	for _, entry := range entries {
		if entry.Name() == name {
			f.currentPathEntry = entry
			return nil
		}
	}

	return fmt.Errorf("Error: Entry %s does not exist in the current directory", name)
}

func (f *FileBrowser) listDir() []os.DirEntry {
	files, err := os.ReadDir(f.currentPath)

	if err != nil {
		log.Printf("Error reading directory %s: %v", f.currentPath, err)
		return []os.DirEntry{}
	}

	return files
}

func (f *FileBrowser) NavigateInto(name string) error {
	entries := f.listDir()
	for _, entry := range entries {
		if entry.Name() == name && entry.IsDir() {
			newPath := filepath.Join(f.currentPath, name)

			if _, err := os.Stat(newPath); err == nil {
				f.currentPath = newPath
				f.currentPathEntry = &dirEntry{
					name:  name,
					isDir: true,
				}
				f.dirEntries = f.listDir()
				return nil
			} else {
				return fmt.Errorf("Error: Cannot access directory %s: %v", name, err)
			}
		}
	}
	return fmt.Errorf("Error: %s is not a valid directory", name)
}

func (f *FileBrowser) NavigateUp() error {
	parentDir := filepath.Dir(f.currentPath)

	if parentDir == f.currentPath {
		return fmt.Errorf("Error: Already at the root directory")
	}

	f.currentPath = parentDir
	f.currentPathEntry = &dirEntry{
		name:  filepath.Base(parentDir),
		isDir: true,
	}
	f.dirEntries = f.listDir()

	return nil
}

func (f *FileBrowser) GetFileContent(file string) (string, error) {
	filePath := filepath.Join(f.currentPath, file)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return "", fmt.Errorf("Error reading file: %v", err)
	}
	return string(content), nil
}

func (f *FileBrowser) SaveFileContent(file string, content string) error {
	err := os.WriteFile(filepath.Join(f.CurrentPath(), file), []byte(content), 0644)

	if err != nil {
		return err
	}

	return nil
}

func (f *FileBrowser) GetFullPath(name string) string {
	return filepath.Join(f.currentPath, name)
}

func (f *FileBrowser) buildSearchIndex() error {
	return nil
}

func (f *FileBrowser) Search(searchInput string) string {
	return searchInput
}
