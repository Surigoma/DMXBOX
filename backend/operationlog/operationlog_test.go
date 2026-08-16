package operationlog

import (
	"path/filepath"
	"testing"
)

func TestOperationLogPersists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "operations.jsonl")
	if err := Open(path); err != nil {
		t.Fatal(err)
	}
	if err := Record("webui", "dmx", "fade", map[string]string{"id": "stage"}); err != nil {
		t.Fatal(err)
	}
	Close()
	if err := Open(path); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(Close)
	entries := List(1)
	if len(entries) != 1 || entries[0].Source != "webui" || entries[0].Args["id"] != "stage" {
		t.Fatalf("unexpected entries: %#v", entries)
	}
}
