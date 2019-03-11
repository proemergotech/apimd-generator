package generator

import (
	"testing"
)

func TestNormalizePath(t *testing.T) {
	for _, data := range []struct {
		Original string
		Want     string
	}{
		{
			Original: "/normal/test/path",
			Want:     "/normal/test/path",
		},
		{
			Original: "/param/:param",
			Want:     "/param/{param}",
		},
		{
			Original: "/param/:param/path\\:escaped",
			Want:     "/param/{param}/path:escaped",
		},
		{
			Original: "/param/:param/path\\:escaped#auth",
			Want:     "/param/{param}/path:escaped#auth",
		},
		{
			Original: "/param/:multiple:param",
			Want:     "/param/{multiple}{param}",
		},
	} {
		got := normalizePath(data.Original)
		if got != data.Want {
			t.Errorf("original: `%s` want: `%s` got: `%s`", data.Original, data.Want, got)
		}
	}
}
