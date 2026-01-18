package sort

import (
	"testing"
	"time"

	"github.com/ipanardian/lu-hut/internal/model"
)

func TestNameSortStrategy(t *testing.T) {
	strategy := &Name{}

	files := []model.FileEntry{
		{Name: "batu-bara.txt", IsDir: false},
		{Name: "sawit.txt", IsDir: false},
		{Name: "bisnis", IsDir: true},
		{Name: "lord.txt", IsDir: false},
	}

	strategy.Sort(files, false)

	expected := []string{"bisnis", "batu-bara.txt", "lord.txt", "sawit.txt"}
	for i, f := range files {
		if f.Name != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, f.Name)
		}
	}
}

func TestTimeSortStrategy(t *testing.T) {
	strategy := &Time{}

	now := time.Now()
	files := []model.FileEntry{
		{Name: "old.txt", ModTime: now.Add(-24 * time.Hour)},
		{Name: "new.txt", ModTime: now.Add(-1 * time.Hour)},
		{Name: "newest.txt", ModTime: now},
	}

	strategy.Sort(files, false)

	expected := []string{"newest.txt", "new.txt", "old.txt"}
	for i, f := range files {
		if f.Name != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, f.Name)
		}
	}
}

func TestSizeSortStrategy(t *testing.T) {
	strategy := &Size{}

	files := []model.FileEntry{
		{Name: "small.txt", Size: 100, IsDir: false},
		{Name: "large.txt", Size: 10000, IsDir: false},
		{Name: "directory", Size: 0, IsDir: true},
		{Name: "medium.txt", Size: 1000, IsDir: false},
	}

	strategy.Sort(files, false)

	expected := []string{"directory", "large.txt", "medium.txt", "small.txt"}
	for i, f := range files {
		if f.Name != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, f.Name)
		}
	}
}

func TestExtensionSortStrategy(t *testing.T) {
	strategy := &Extension{}

	files := []model.FileEntry{
		{Name: "file.txt", IsDir: false},
		{Name: "file.go", IsDir: false},
		{Name: "directory", IsDir: true},
		{Name: "file.py", IsDir: false},
	}

	strategy.Sort(files, false)

	expected := []string{"directory", "file.go", "file.py", "file.txt"}
	for i, f := range files {
		if f.Name != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, f.Name)
		}
	}
}
