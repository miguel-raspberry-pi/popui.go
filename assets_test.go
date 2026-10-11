package popui

import (
	"context"
	"strings"
	"testing"
)

func TestAssetURL(t *testing.T) {
	ctx := context.Background()
	if got := AssetURL(ctx, "assets", "popui.css"); !strings.HasPrefix(got, AssetPath+"/assets/popui.css?v=") {
		t.Errorf("default base: got %q", got)
	}
	ctx = WithAssetBase(ctx, "/apps/demo"+AssetPath)
	if got := AssetURL(ctx, "assets", "popui.css"); !strings.HasPrefix(got, "/apps/demo/_popui/assets/popui.css?v=") {
		t.Errorf("custom base: got %q", got)
	}
}
