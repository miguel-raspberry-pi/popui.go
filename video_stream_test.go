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
	for _, unwanted := range []string{`data-video-stream-controls`, `data-ice-servers`} {
		if strings.Contains(manual, unwanted) {
			t.Errorf("unexpected %s in %s", unwanted, manual)
		}
	}
}
