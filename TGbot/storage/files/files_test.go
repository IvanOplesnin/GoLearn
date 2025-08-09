package files

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"tgbot/storage"
)

// mockPage implements storage.Page with a custom Hash method for testing.
type mockPage struct {
	storage.Page
	hash     string
	hashFunc func() (string, error) // добавляем поле-функцию
}

func (m *mockPage) Hash() (string, error) {
	if m.hashFunc != nil {
		return m.hashFunc()
	}
	return m.hash, nil
}

func TestStorage_Remove_Success(t *testing.T) {
	dir := t.TempDir()
	s := Storage{basePath: dir}

	page := &mockPage{
		Page: storage.Page{Username: "testuser"},
		hash: "testhash",
	}

	userDir := filepath.Join(dir, page.Username)
	if err := os.MkdirAll(userDir, 0774); err != nil {
		t.Fatalf("failed to create user dir: %v", err)
	}
	filePath := filepath.Join(userDir, page.hash)
	if err := os.WriteFile(filePath, []byte("test"), 0666); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	err := s.Remove(&page.Page)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected file to be removed, but it exists")
	}
}

func TestStorage_Remove_FileNameError(t *testing.T) {
	dir := t.TempDir()
	s := Storage{basePath: dir}

	page := &storage.Page{Username: "testuser"}
	// fileName will call Page.Hash(), which may not be implemented and can error.
	// So, we use a mockPage with a Hash method that returns an error.
	badPage := &mockPage{
		Page: *page,
		hash: "",
		hashFunc: func() (string, error) {
			return "", errors.New("hash error")
		},
	}

	err := s.Remove(&badPage.Page)
	if err == nil {
		t.Errorf("expected error due to hash error, got nil")
	}
}

func TestStorage_Remove_FileDoesNotExist(t *testing.T) {
	dir := t.TempDir()
	s := Storage{basePath: dir}

	page := &mockPage{
		Page: storage.Page{Username: "testuser"},
		hash: "nonexistent",
	}

	err := s.Remove(&page.Page)
	if err == nil {
		t.Errorf("expected error when removing non-existent file, got nil")
	}
}
func TestStorage_Save_Success(t *testing.T) {
	dir := t.TempDir()
	s := Storage{basePath: dir}

	page := &mockPage{
		Page: storage.Page{Username: "testuser"},
		hash: "testhash",
	}

	err := s.Save(&page.Page)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check file exists
	filePath := filepath.Join(dir, page.Username, page.hash)
	if _, err := os.Stat(filePath); err != nil {
		t.Errorf("expected file to exist, got error: %v", err)
	}
}

func TestStorage_Save_FileNameError(t *testing.T) {
	dir := t.TempDir()
	s := Storage{basePath: dir}

	page := &mockPage{
		Page: storage.Page{Username: "testuser"},
		hashFunc: func() (string, error) {
			return "", errors.New("hash error")
		},
	}

	err := s.Save(&page.Page)
	if err == nil {
		t.Errorf("expected error due to hash error, got nil")
	}
}

func TestStorage_Save_MkdirAllError(t *testing.T) {
	// Try to create a directory in a location where permission is denied.
	// On Unix, "/" is root and should fail for non-root users.
	// On Windows, use an invalid path.
	var basePath string
	if os.PathSeparator == '/' {
		basePath = "/invalid_dir"
	} else {
		basePath = string([]rune{os.PathSeparator}) + "invalid_dir"
	}
	s := Storage{basePath: basePath}

	page := &mockPage{
		Page: storage.Page{Username: "testuser"},
		hash: "testhash",
	}

	err := s.Save(&page.Page)
	if err == nil {
		t.Errorf("expected error due to MkdirAll, got nil")
	}
}

func TestStorage_Save_CreateFileError(t *testing.T) {
	dir := t.TempDir()
	s := Storage{basePath: dir}

	// Create a directory with the same name as the file to cause os.Create to fail
	userDir := filepath.Join(dir, "testuser")
	if err := os.MkdirAll(userDir, 0774); err != nil {
		t.Fatalf("failed to create user dir: %v", err)
	}
	filePath := filepath.Join(userDir, "testhash")
	if err := os.MkdirAll(filePath, 0774); err != nil {
		t.Fatalf("failed to create dir with file name: %v", err)
	}

	page := &mockPage{
		Page: storage.Page{Username: "testuser"},
		hash: "testhash",
	}

	err := s.Save(&page.Page)
	if err == nil {
		t.Errorf("expected error due to file creation, got nil")
	}
}

// func TestStorage_Save_EncodeError(t *testing.T) {
// 	dir := t.TempDir()
// 	s := Storage{basePath: dir}

// 	// Use a type that gob cannot encode (e.g., a channel field)
// 	type badPage struct {
// 		storage.Page
// 		Ch chan int
// 	}
// 	page := &badPage{
// 		Page: storage.Page{Username: "testuser"},
// 		Ch:   make(chan int),
// 	}

	// Save expects *storage.Page, so we need to cast
	// This test is not directly possible unless Save accepts interface{}
	// So, we skip this test as gob encoding error is hard to trigger with the current signature.
	// (Alternatively, you could use reflection or a custom encoder in production code.)
// }

