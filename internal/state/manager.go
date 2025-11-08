package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type KeyExtractor[K comparable, S any] func(*S) (K, bool)
type Initializer[S any] func(*S)

type Manager[K comparable, S any] struct {
	defaultAppName string
	stateFileName  string
	stateFilePath  string
	states         map[K]*S
	keyExtractor   KeyExtractor[K, S]
	initializer    Initializer[S]
}

func NewManager[K comparable, S any](defaultAppName string, keyExtractor KeyExtractor[K, S], initializer Initializer[S]) *Manager[K, S] {
	m := &Manager[K, S]{defaultAppName: defaultAppName, keyExtractor: keyExtractor, initializer: initializer, states: make(map[K]*S), stateFileName: buildStateFileName(defaultAppName)}
	m.stateFilePath = m.computeStateFilePath()
	_ = m.loadStates()
	return m
}

func (m *Manager[K, S]) computeStateFilePath() string {
	// Prefer HOME environment variable if set so tests using t.Setenv("HOME", ...) work across platforms.
	// On Windows, os.UserHomeDir may ignore a test-set HOME, leading to leakage of existing states and test flakiness.
	home := os.Getenv("HOME")
	if strings.TrimSpace(home) == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil || strings.TrimSpace(home) == "" {
			home = "."
		}
	}
	fileName := m.stateFileName
	if fileName == "" {
		fileName = buildStateFileName(m.defaultAppName)
	}
	return filepath.Join(home, fileName)
}

func buildStateFileName(appName string) string {
	name := strings.TrimSpace(appName)
	if name == "" {
		return ""
	}
	name = strings.TrimPrefix(name, ".")
	name = strings.TrimSuffix(name, ".json")
	name = strings.ToLower(name)
	return fmt.Sprintf(".%s.json", name)
}

func (m *Manager[K, S]) Configure(appName string) error {
	if appName == "" {
		appName = m.defaultAppName
	}
	m.stateFileName = buildStateFileName(appName)
	if m.stateFileName == "" {
		return fmt.Errorf("state: invalid application name")
	}
	m.stateFilePath = m.computeStateFilePath()
	m.states = make(map[K]*S)
	return m.loadStates()
}

func (m *Manager[K, S]) Get(key K) *S { return m.states[key] }

func (m *Manager[K, S]) Set(state S) error {
	if m.keyExtractor == nil {
		return fmt.Errorf("state: key extractor is not defined")
	}
	stateCopy := state
	key, ok := m.keyExtractor(&stateCopy)
	if !ok {
		return fmt.Errorf("state: could not extract key")
	}
	m.states[key] = &stateCopy
	return m.saveStates()
}

func (m *Manager[K, S]) Delete(key K) error { delete(m.states, key); return m.saveStates() }
func (m *Manager[K, S]) Save() error        { return m.saveStates() }

func (m *Manager[K, S]) AllStates() map[K]*S {
	copyMap := make(map[K]*S, len(m.states))
	for k, v := range m.states {
		copyMap[k] = v
	}
	return copyMap
}

// WipeFile deletes the backing state file and clears in-memory states.
// If the file does not exist, it silently succeeds.
func (m *Manager[K, S]) WipeFile() error {
	if m.stateFilePath == "" {
		return fmt.Errorf("state: state file path is not set")
	}
	if err := os.Remove(m.stateFilePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	m.states = make(map[K]*S)
	return nil
}

func (m *Manager[K, S]) loadStates() error {
	if m.stateFilePath == "" {
		return fmt.Errorf("state: state file path is not set")
	}
	data, err := os.ReadFile(m.stateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var states []S
	if err := json.Unmarshal(data, &states); err != nil {
		return err
	}
	for i := range states {
		entry := states[i]
		if m.initializer != nil {
			m.initializer(&entry)
		}
		if m.keyExtractor == nil {
			return fmt.Errorf("state: key extractor is not defined")
		}
		key, ok := m.keyExtractor(&entry)
		if !ok {
			continue
		}
		entryCopy := entry
		m.states[key] = &entryCopy
	}
	return nil
}

func (m *Manager[K, S]) saveStates() error {
	if m.stateFilePath == "" {
		return fmt.Errorf("state: state file path is not set")
	}
	states := make([]S, 0, len(m.states))
	for _, state := range m.states {
		states = append(states, *state)
	}
	data, err := json.MarshalIndent(states, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.stateFilePath, data, 0644)
}

func BuildStateFileName(appName string) (string, error) {
	name := strings.TrimSpace(appName)
	if name == "" {
		return "", fmt.Errorf("state: application name cannot be empty")
	}
	return buildStateFileName(name), nil
}

func (m *Manager[K, S]) StateFilePath() string { return m.stateFilePath }
