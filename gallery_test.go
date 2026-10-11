package popui

import (
	"strings"
	"testing"

	"github.com/miguel-raspberry-pi/popui.go/props"
)

func TestGallery(t *testing.T) {
	got := renderString(t, Gallery(props.Gallery{
		ID: "photos",
		Images: []props.GalleryImage{
			{Src: "/img/1.jpg", Thumb: "/img/1-thumb.jpg", Alt: "Living room", Caption: "Bright living room"},
			{Src: "/img/2.jpg", Alt: "Kitchen <new>"},
			{Src: "/img/3.jpg"},
		},
	}))
	for _, want := range []string{
		`id="photos"`,
		`data-gallery`,
		`x-data="{ i: 0, imgs: [{&#34;src&#34;:&#34;/img/1.jpg&#34;`,
		`\u003cnew\u003e`, // JSON escapes < and > inside the Alpine scope
		`src="/img/1.jpg" alt="Living room"`,
		`href="/img/1.jpg"`,
		`@click.prevent="show()"`,
		`@keydown.arrow-left.prevent="go(-1)"`,
		`>1 / 3</span>`,
		`>Bright living room</p>`,
		`src="/img/1-thumb.jpg"`,
		`src="/img/2.jpg" alt="" loading="lazy"`,
		`data-state="active" aria-current="true"`,
		`data-state="inactive"`,
		`@click="i = 2"`,
		`aria-label="Show photo 2 / 3: Kitchen &lt;new&gt;"`,
		`<dialog x-ref="viewer"`,
		`aria-label="Previous photo"`,
		`aria-label="Close viewer"`,
		`x-cloak`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in %s", want, got)
		}
	}
	// Main image box plus three thumbnails.
	if n := strings.Count(got, "aspect-[4/3]"); n != 4 {
		t.Errorf("aspect-[4/3] used %d times", n)
	}
	if strings.Contains(got, ">else<") || strings.Contains(got, " else ") {
		t.Errorf("stray else in %s", got)
	}

	single := renderString(t, Gallery(props.Gallery{AspectClass: "aspect-video", Images: []props.GalleryImage{{Src: "/a.jpg"}}}))
	for _, unwanted := range []string{`data-gallery-thumbs`, `data-gallery-counter`, `data-gallery-caption`, `@keydown`, `aspect-[4/3]`} {
		if strings.Contains(single, unwanted) {
			t.Errorf("unexpected %s in %s", unwanted, single)
		}
	}
	if !strings.Contains(single, "aspect-video") {
		t.Errorf("missing aspect class in %s", single)
	}

	empty := renderString(t, Gallery(props.Gallery{Empty: "No photos yet"}))
	if !strings.Contains(empty, "No photos yet") || strings.Contains(empty, "x-data") {
		t.Errorf("empty = %s", empty)
	}
	if def := renderString(t, Gallery()); !strings.Contains(def, "No photos") {
		t.Errorf("default empty = %s", def)
	}
}
