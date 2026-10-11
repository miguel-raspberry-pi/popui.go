package popui

import (
	"strings"
	"testing"

	"github.com/miguel-raspberry-pi/popui.go/props"
)

func TestVideoStream(t *testing.T) {
	got := renderString(t, VideoStream(props.VideoStream{
		WHEPURL:    "/apps/camera/whep?cam=1&x=2",
		Title:      "Front door",
		Controls:   true,
		ICEServers: `[{"urls":"stun:stun.example.org"}]`,
		ControlURL: "/apps/camera/api/stream",
	}))
	for _, want := range []string{
		`data-video-stream`,
		`data-whep-url="/apps/camera/whep?cam=1&amp;x=2"`,
		`data-state="connecting"`,
		`role="region"`,
		`aria-label="Front door"`,
		`style="aspect-ratio: 16/9;"`,
		`data-ice-servers="[{&#34;urls&#34;:&#34;stun:stun.example.org&#34;}]"`,
		`playsinline`,
		`muted`,
		`autoplay`,
		`data-video-stream-overlay="connecting"`,
		`data-video-stream-overlay="live"`,
		`data-video-stream-overlay="offline"`,
		`data-video-stream-retry`,
		`data-video-stream-toggle`,
		`data-video-stream-fullscreen`,
		`data-control-url="/apps/camera/api/stream"`,
		`data-video-stream-overlay="stopped"`,
		`data-video-stream-start`,
		`data-video-stream-live`,
		`data-video-stream-live-label="stop"`,
		`data-video-stream-live-label="start"`,
		`Stop live`,
		`Start live`,
		`data-video-stream-snapshot`,
		`aria-label="Take snapshot"`,
		`data-video-stream-toast`,
		`Snapshot saved`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in %s", want, got)
		}
	}
	if strings.Contains(got, "data-manual") {
		t.Error("data-manual set without Manual")
	}

	manual := renderString(t, VideoStream(props.VideoStream{WHEPURL: "/whep", Manual: true, AspectRatio: "4/3"}))
	for _, want := range []string{`data-manual`, `data-state="idle"`, `aria-label="Live video"`, `style="aspect-ratio: 4/3;"`} {
		if !strings.Contains(manual, want) {
			t.Errorf("missing %s in %s", want, manual)
		}
	}
	for _, unwanted := range []string{`data-video-stream-controls`, `data-ice-servers`, `data-control-url`, `data-video-stream-overlay="stopped"`, `data-video-stream-snapshot`, `data-video-stream-toast`} {
		if strings.Contains(manual, unwanted) {
			t.Errorf("unexpected %s in %s", unwanted, manual)
		}
	}
	// The snapshot button starts disabled: popui.js enables it while live.
	if i := strings.Index(got, `data-video-stream-snapshot`); i < 0 || !strings.Contains(got[strings.LastIndex(got[:i], "<button"):i], "disabled") {
		t.Errorf("snapshot button not disabled initially")
	}
}

func TestVideoStreamControlVariants(t *testing.T) {
	// The live switch alone shows the controls bar with just its button.
	switchOnly := renderString(t, VideoStream(props.VideoStream{WHEPURL: "/whep", ControlURL: "/api/stream"}))
	for _, want := range []string{`data-video-stream-controls`, `data-video-stream-live`, `data-video-stream-overlay="stopped"`} {
		if !strings.Contains(switchOnly, want) {
			t.Errorf("ControlURL only: missing %s", want)
		}
	}
	for _, unwanted := range []string{`data-video-stream-toggle`, `data-video-stream-snapshot`, `data-video-stream-fullscreen`} {
		if strings.Contains(switchOnly, unwanted) {
			t.Errorf("ControlURL only: unexpected %s", unwanted)
		}
	}

	// NoSnapshot drops the button and its toast.
	noSnap := renderString(t, VideoStream(props.VideoStream{WHEPURL: "/whep", Controls: true, NoSnapshot: true}))
	for _, unwanted := range []string{`data-video-stream-snapshot`, `data-video-stream-toast`, `data-video-stream-live`, `data-control-url`} {
		if strings.Contains(noSnap, unwanted) {
			t.Errorf("NoSnapshot: unexpected %s", unwanted)
		}
	}
	if !strings.Contains(noSnap, `data-video-stream-toggle`) {
		t.Error("NoSnapshot: missing pause toggle")
	}
}
