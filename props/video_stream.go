package props

import "github.com/a-h/templ"

// VideoStream Templ component props.
//
// The player connects to a WHEP-style signalling endpoint (non-trickle):
// popui.js POSTs an SDP offer to WHEPURL with Content-Type application/sdp
// and expects 201 Created with the SDP answer as the body and a Location
// header naming the session, which is DELETEd when the player stops. A 503
// response means the stream is offline; the player retries with backoff.
//
// ControlURL adds a server-side live switch: GET ControlURL returns JSON
// with a boolean "enabled" (other fields are ignored), and POST
// ControlURL+"/stop" and ControlURL+"/start" change it, answering with the
// same JSON. While it is off the player shows data-state="stopped" instead
// of retrying.
type VideoStream struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// WHEPURL is the signalling endpoint, e.g. "/apps/camera/whep". Required.
	WHEPURL string
	// Title labels the video for assistive technology, e.g. "Front door
	// camera". Defaults to "Live video".
	Title string
	// Manual waits for the viewer to press Play instead of connecting as
	// soon as the player is on the page. The zero value connects on load.
	Manual bool
	// Controls shows snapshot, pause/play and fullscreen buttons over the
	// video. Pause only ends this viewer's session.
	Controls bool
	// NoSnapshot hides the Snapshot button that Controls otherwise shows. It
	// saves the current frame client-side as camera-YYYYMMDD-HHMMSS.png.
	NoSnapshot bool
	// ControlURL is the server-side live switch (see above), e.g.
	// "/apps/camera/api/stream". It shows a Stop live / Start live button
	// that turns the stream off for every viewer, and the player polls it
	// every 5s while the page is visible. Optional.
	ControlURL string
	// AspectRatio is the CSS aspect-ratio of the player box, e.g. "4/3".
	// Defaults to "16/9".
	AspectRatio string
	// ICEServers is an optional JSON array of RTCIceServer objects, e.g.
	// `[{"urls":"stun:stun.example.org"}]`. Not needed on a LAN.
	ICEServers string
}
