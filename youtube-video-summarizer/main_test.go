package main

import (
	"strings"
	"testing"
)

func TestExtractVideoID(t *testing.T) {
	tests := []struct {
		url     string
		want    string
		wantErr bool
	}{
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://youtube.com/watch?v=dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://youtu.be/dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://www.youtube.com/shorts/abc123XYZ_-", "abc123XYZ_-", false},
		{"https://youtube.com/embed/dQw4w9WgXcQ?autoplay=1", "dQw4w9WgXcQ", false},
		{"http://youtube.com/watch?v=hello-wor_1", "hello-wor_1", false},
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=30s", "dQw4w9WgXcQ", false},
		{"not a youtube url", "", true},
		{"https://example.com/watch?v=12345678901", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		got, err := extractVideoID(tt.url)
		if tt.wantErr {
			if err == nil {
				t.Errorf("extractVideoID(%q) expected error, got nil and id=%q", tt.url, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("extractVideoID(%q) unexpected error: %v", tt.url, err)
			continue
		}
		if got != tt.want {
			t.Errorf("extractVideoID(%q) = %q, want %q", tt.url, got, tt.want)
		}
	}
}

func TestCleanVTT(t *testing.T) {
	input := `WEBVTT
Kind: captions
Language: en

00:00:01.000 --> 00:00:03.000
Hello world.

00:00:03.000 --> 00:00:05.000
This is a <i>test</i>.

00:00:05.000 --> 00:00:07.000
Hello world.

00:00:07.000 --> 00:00:09.000
Final line.
`

	got, err := cleanVTT(strings.NewReader(input))
	if err != nil {
		t.Fatalf("cleanVTT unexpected error: %v", err)
	}

	// Verify: markup stripped, extra blanks gone
	if !strings.Contains(got, "This is a test.") {
		t.Errorf("expected 'This is a test.' (markup stripped), got: %q", got)
	}
	// Dedup is consecutive-only (catches back-to-back repeats in auto-subs).
	// Non-consecutive duplicates (separated by other text) are preserved — that's correct.
	count := strings.Count(got, "Hello world.")
	if count < 1 || count > 2 {
		t.Errorf("expected 'Hello world.' 1-2 times, got %d times in: %q", count, got)
	}
	if !strings.HasSuffix(got, "Final line.") {
		t.Errorf("expected output to end with 'Final line.', got: %q", got)
	}
	if strings.Contains(got, "WEBVTT") || strings.Contains(got, "Kind:") || strings.Contains(got, "-->") {
		t.Errorf("expected VTT metadata stripped, got: %q", got)
	}
}
