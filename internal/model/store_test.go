package model

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/adrg/xdg"
)

func TestModelPersistence(t *testing.T) {
	tempDir := t.TempDir()
	dataDir := filepath.Join(tempDir, "data")
	configDir := filepath.Join(tempDir, "config")

	t.Setenv("XDG_DATA_HOME", dataDir)
	t.Setenv("XDG_CONFIG_HOME", configDir)
	t.Setenv("LOCALAPPDATA", dataDir)
	t.Setenv("APPDATA", configDir)
	xdg.Reload()
	t.Cleanup(func() {
		xdg.Reload()
	})

	// Test tasks
	tasks, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}

	testTasks := []Task{
		{
			ID:        "t-1",
			Title:     "Buy milk",
			Done:      false,
			Date:      "2026-09-07",
			CreatedAt: time.Now(),
		},
	}
	if err := Save(testTasks); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	loadedTasks, err := Load()
	if err != nil {
		t.Fatalf("Load() after save failed: %v", err)
	}
	if len(loadedTasks) != 1 || loadedTasks[0].Title != "Buy milk" {
		t.Fatalf("unexpected loaded tasks: %+v", loadedTasks)
	}

	// Test notes
	notes, err := LoadNotes()
	if err != nil {
		t.Fatalf("LoadNotes() failed: %v", err)
	}
	if len(notes) != 0 {
		t.Fatalf("expected 0 notes, got %d", len(notes))
	}

	testNotes := []Note{
		{
			ID:        "n-1",
			Body:      "Remember to call Alice",
			CreatedAt: time.Now(),
		},
	}
	if err := SaveNotes(testNotes); err != nil {
		t.Fatalf("SaveNotes() failed: %v", err)
	}

	loadedNotes, err := LoadNotes()
	if err != nil {
		t.Fatalf("LoadNotes() after save failed: %v", err)
	}
	if len(loadedNotes) != 1 || loadedNotes[0].Body != "Remember to call Alice" {
		t.Fatalf("unexpected loaded notes: %+v", loadedNotes)
	}

	// Test settings
	settings, err := LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings() failed: %v", err)
	}
	if settings.Theme != "" {
		t.Fatalf("expected empty default settings, got %+v", settings)
	}

	testSettings := Settings{
		Theme:                     "dracula",
		Layout:                    "standard",
		PomodoroFocusMinutes:      25,
		PomodoroShortBreakMinutes: 5,
	}
	if err := SaveSettings(testSettings); err != nil {
		t.Fatalf("SaveSettings() failed: %v", err)
	}

	loadedSettings, err := LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings() after save failed: %v", err)
	}
	if loadedSettings.Theme != "dracula" || loadedSettings.PomodoroFocusMinutes != 25 {
		t.Fatalf("unexpected loaded settings: %+v", loadedSettings)
	}

	// Verify the actual files exist in the XDG directories
	tasksPath, err := dataPath()
	if err != nil {
		t.Fatalf("dataPath() failed: %v", err)
	}
	if _, err := os.Stat(tasksPath); err != nil {
		t.Fatalf("expected tasks file at %s, stat error: %v", tasksPath, err)
	}

	notesFilePath, err := notesPath()
	if err != nil {
		t.Fatalf("notesPath() failed: %v", err)
	}
	if _, err := os.Stat(notesFilePath); err != nil {
		t.Fatalf("expected notes file at %s, stat error: %v", notesFilePath, err)
	}

	settingsFilePath, err := settingsPath()
	if err != nil {
		t.Fatalf("settingsPath() failed: %v", err)
	}
	if _, err := os.Stat(settingsFilePath); err != nil {
		t.Fatalf("expected settings file at %s, stat error: %v", settingsFilePath, err)
	}
}
