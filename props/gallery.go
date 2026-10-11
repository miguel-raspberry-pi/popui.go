package props

import "github.com/a-h/templ"

// Gallery defines the props of a photo gallery: a large main image, a strip
// of thumbnails that switch it, and a fullscreen viewer opened from the main
// image. Interaction uses Alpine.js (bundled by popui.App); without
// JavaScript the first image and its caption are shown.
type Gallery struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Images to show, in order. The first one is shown initially.
	Images []GalleryImage
	// Empty is the title of the empty state shown when there are no
	// images. Defaults to "No photos".
	Empty string
	// AspectClass is the Tailwind aspect-ratio class of the main image box,
	// e.g. "aspect-video". Defaults to "aspect-[4/3]". Images are cropped
	// to fill it; the fullscreen viewer shows them whole.
	AspectClass string
}

// GalleryImage is one photo of a Gallery.
type GalleryImage struct {
	// Src is the URL of the full-size image. Required.
	Src string
	// Thumb is the URL of a smaller version for the thumbnail strip.
	// Defaults to Src.
	Thumb string
	// Alt describes the image for assistive technology.
	Alt string
	// Caption is shown under the main image and in the viewer.
	Caption string
}
