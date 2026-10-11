package props

import "github.com/a-h/templ"

// Frame Templ component props.
type Frame struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Src is the page to embed. Prefer same-origin paths: pages built with
	// PopUI usually only allow being framed by their own origin.
	Src templ.SafeURL
	// Title describes the embedded page for assistive technology.
	Title string
	// Lazy defers loading until the frame is about to become visible.
	Lazy bool
}
