package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseExcludes(t *testing.T) {
	home := os.Getenv("HOME")
	startDir := "/jd"

	cases := []struct {
		name string
		raw  string
		want []string
	}{
		{"empty", "", nil},
		{"whitespace only", "   ,  ,", nil},
		{"relative entries joined to startDir", "99_archive, 30_drafts",
			[]string{"/jd/99_archive", "/jd/30_drafts"}},
		{"absolute path preserved", "/abs/path", []string{"/abs/path"}},
		{"tilde expanded", "~/foo", []string{filepath.Join(home, "foo")}},
		{"trailing slash cleaned", "99_archive/", []string{"/jd/99_archive"}},
		{"mixed entries", "  99_archive , ~/foo , /etc ",
			[]string{"/jd/99_archive", filepath.Join(home, "foo"), "/etc"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseExcludes(tc.raw, startDir)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parseExcludes(%q) = %v, want %v", tc.raw, got, tc.want)
			}
		})
	}
}

func TestIsExcluded(t *testing.T) {
	excludes := []string{"/jd/99_archive", "/jd/30_drafts"}

	cases := []struct {
		path string
		want bool
	}{
		{"/jd/99_archive", true},
		{"/jd/99_archive/12.34 something", true},
		{"/jd/30_drafts/sub/deep", true},
		{"/jd/10_active", false},
		{"/jd/99_archive_other", false}, // prefix only, not a path boundary
		{"", false},
	}

	for _, tc := range cases {
		got := isExcluded(tc.path, excludes)
		if got != tc.want {
			t.Errorf("isExcluded(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}

	if isExcluded("/jd/anything", nil) {
		t.Error("isExcluded with nil excludes should be false")
	}
}
