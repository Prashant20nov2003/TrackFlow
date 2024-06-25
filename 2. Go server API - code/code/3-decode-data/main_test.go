package main

import (
	"testing"
)

func TestDecodeData(t *testing.T) {
	data, err := decodeData("eyJ0cmFja2luZyI6eyJ0eXBlIjoicGFnZSIsImlkZW50aXR5IjoiIiwidWEiOiJNb3ppbGxhLzUuMCAoV2luZG93cyBOVCAxMC4wOyBXaW42NDsgeDY0KSBBcHBsZVdlYktpdC81MzcuMzYgKEtIVE1MLCBsaWtlIEdlY2tvKSBDaHJvbWUvMTE5LjAuMC4wIFNhZmFyaS81MzcuMzYiLCJldmVudCI6Ii8iLCJjYXRlZ29yeSI6IlBhZ2Ugdmlld3MiLCJyZWZlcnJlciI6IiIsImlzVG91Y2hEZXZpY2UiOmZhbHNlfSwic2l0ZV9pZCI6Im15LXNpdGUtaWQtaGVyZSJ9")
	if err != nil {
		t.Fatal(err)
	} else if data.SiteID != "my-site-id-here" {
		t.Errorf("expected 'my-site-id-here' got %s", data.SiteID)
	}
}
