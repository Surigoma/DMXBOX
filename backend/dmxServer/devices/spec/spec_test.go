package spec_test

import (
	"backend/dmxServer/devices/spec"
	"testing"
)

func Test_NewSpec(t *testing.T) {
	t.Run("WCLight", func(t *testing.T) {
		got := spec.NewWCLight()
		if got == nil {
			t.Error("Failed to create WCLight")
			return
		}
		if got.UseChannel != 3 {
			t.Errorf("Parameter error. UseChannel: %v want 3", got.UseChannel)
		}
		if got.Model != "wclight" {
			t.Errorf("Parameter error. Name: %v want wclight", got.Model)
		}
	})
	t.Run("Dimmer", func(t *testing.T) {
		got := spec.NewDimmer()
		if got == nil {
			t.Error("Failed to create Dimmer")
			return
		}
		if got.UseChannel != 1 {
			t.Errorf("Parameter error. UseChannel: %v want 1", got.UseChannel)
		}
		if got.Model != "dimmer" {
			t.Errorf("Parameter error. Name: %v want dimmer", got.Model)
		}
	})
}
