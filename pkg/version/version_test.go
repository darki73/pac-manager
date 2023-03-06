package version

import (
	"testing"
	"time"
)

func TestGetVersion(t *testing.T) {
	got := GetVersion()
	want := "dev"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetCommit(t *testing.T) {
	got := GetCommit()
	want := "none"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetBuildDate(t *testing.T) {
	got := GetBuildDate()
	want := "0000-00-00T00-00-00Z"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetBuiltBy(t *testing.T) {
	got := GetBuiltBy()
	want := "freedomcore"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestGetStartDateTime(t *testing.T) {
	got := GetStartDateTime()
	want := time.Now()
	before := want.Add(time.Duration(-5) * time.Second)
	after := want.Add(time.Duration(5) * time.Second)

	if !inTimeSpan(before, after, got) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}
