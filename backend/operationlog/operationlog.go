package operationlog

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
	"time"
)

const maxEntries = 1000

// ponytail: the append-only file has no rotation; add size-based rotation if operation volume makes it necessary.

type Entry struct {
	Time   time.Time         `json:"time"`
	Source string            `json:"source,omitempty"`
	Target string            `json:"target"`
	Action string            `json:"action"`
	Args   map[string]string `json:"args"`
}

var store struct {
	sync.RWMutex
	entries []Entry
	file    *os.File
}

func Open(path string) error {
	store.Lock()
	defer store.Unlock()
	if store.file != nil {
		_ = store.file.Close()
	}
	store.entries = nil
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		store.file = nil
		return err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry Entry
		if json.Unmarshal(scanner.Bytes(), &entry) == nil {
			store.entries = append(store.entries, entry)
			if len(store.entries) > maxEntries {
				store.entries = store.entries[1:]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		_ = file.Close()
		store.file = nil
		return err
	}
	store.file = file
	return nil
}

func Record(source, target, action string, args map[string]string) error {
	entry := Entry{Time: time.Now(), Source: source, Target: target, Action: action, Args: make(map[string]string, len(args))}
	for key, value := range args {
		entry.Args[key] = value
	}
	store.Lock()
	defer store.Unlock()
	store.entries = append(store.entries, entry)
	if len(store.entries) > maxEntries {
		store.entries = store.entries[1:]
	}
	if store.file != nil {
		return json.NewEncoder(store.file).Encode(entry)
	}
	return nil
}

func List(limit int) []Entry {
	store.RLock()
	defer store.RUnlock()
	if limit <= 0 || limit > len(store.entries) {
		limit = len(store.entries)
	}
	result := make([]Entry, limit)
	copy(result, store.entries[len(store.entries)-limit:])
	return result
}

func Close() {
	store.Lock()
	defer store.Unlock()
	if store.file != nil {
		_ = store.file.Close()
		store.file = nil
	}
}
