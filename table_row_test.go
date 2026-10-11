package popui

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/miguel-raspberry-pi/popui.go/props"
)

func renderString(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestTableRow(t *testing.T) {
	got := renderString(t, TableRow(props.TableRow{Href: "/listings/1?a=1&b=2", Target: "_blank"}))
	for _, want := range []string{`data-href="/listings/1?a=1&amp;b=2"`, `data-target="_blank"`, `tabindex="0"`, `class="popui-table-row-link"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in %s", want, got)
		}
	}

	plain := renderString(t, TableRow())
	if plain != "<tr></tr>" {
		t.Errorf("plain row = %s", plain)
	}
}
