package docs_test

import (
	"testing"

	"github.com/swaggo/swag"
)

func TestDocs(t *testing.T) {
	t.Run("Docs test", func(t *testing.T) {
		result, err := swag.ReadDoc()
		if err != nil {
			t.Error("Fail", "err", err)
		}
		if result == "" {
			t.Error("Failed to generate docs")
		}
		t.Log("Can generate OpenAPI: length", len(result))
	})
}
