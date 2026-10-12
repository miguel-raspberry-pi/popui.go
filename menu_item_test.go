package popui

import (
	"context"
	"strings"
	"testing"

	"github.com/miguel-raspberry-pi/popui.go/props"
)

func renderMenuItem(t *testing.T, p props.MenuItem) string {
	t.Helper()
	var b strings.Builder
	if err := MenuItem(p).Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestMenuItemActions(t *testing.T) {
	for name, tc := range map[string]struct {
		p          props.MenuItem
		want, deny []string
	}{
		"plain":         {props.MenuItem{}, []string{"<div class="}, []string{"<a ", "<button"}},
		"link":          {props.MenuItem{Href: "/x", Target: "_blank"}, []string{`<a href="/x" target="_blank" role="menuitem"`}, []string{"<button"}},
		"submit":        {props.MenuItem{Type: "submit", Form: "f", Name: "enabled", Value: "on"}, []string{`<button type="submit" role="menuitem" form="f" name="enabled" value="on"`}, []string{"<a "}},
		"popover":       {props.MenuItem{PopoverTarget: "confirm"}, []string{`type="button"`, `popovertarget="confirm"`}, nil},
		"disabled link": {props.MenuItem{Href: "/x", Disabled: true}, []string{"pointer-events-none", "disabled"}, []string{"<a "}},
	} {
		html := renderMenuItem(t, tc.p)
		for _, w := range tc.want {
			if !strings.Contains(html, w) {
				t.Errorf("%s: missing %q in %s", name, w, html)
			}
		}
		for _, d := range tc.deny {
			if strings.Contains(html, d) {
				t.Errorf("%s: unexpected %q in %s", name, d, html)
			}
		}
	}
}
