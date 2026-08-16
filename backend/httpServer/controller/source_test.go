package controller

import (
	"backend/config"
	"backend/message"
	"backend/operationlog"
	"backend/packageModule"
	"log/slog"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSendControlUsesPackageModuleName(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(t.Output(), nil))
	if err := operationlog.Open(filepath.Join(t.TempDir(), "operations.jsonl")); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(operationlog.Close)
	manager := packageModule.GetModuleManager()
	manager.Initialize(logger)
	target := packageModule.PackageModule{
		ModuleName: "target",
		Initialize: func(*packageModule.PackageModule, *config.Config) bool { return true },
		Run:        func() {}, Stop: func() {},
		MessageHandler: func(message.Message) int { return 0 },
	}
	manager.RegisterModule("target", &target)
	manager.ModuleInitialize(logger, "test")
	source := packageModule.PackageModule{ModuleName: "renamed-http-module"}
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set(controlModuleKey, &source)
	if !SendControl(ctx, message.Message{To: "target", Arg: message.MessageBody{Action: "fade"}}) {
		t.Fatal("SendControl() failed")
	}
	entries := operationlog.List(1)
	if len(entries) != 1 || entries[0].Source != source.ModuleName {
		t.Fatalf("unexpected entries: %#v", entries)
	}
}
