package initialization

import (
	"strings"
	"testing"
)

func TestBrowserCommand(t *testing.T) {
	tests := []struct {
		name      string
		goos      string
		want      string
		wantFirst string
	}{
		{name: "windows", goos: "windows", want: "rundll32", wantFirst: "url.dll"},
		{name: "darwin", goos: "darwin", want: "open", wantFirst: "http://127.0.0.1:8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, args, err := browserCommand(tt.goos, "http://127.0.0.1:8080")
			if err != nil {
				t.Fatalf("browserCommand() error = %v", err)
			}
			if command != tt.want {
				t.Fatalf("browserCommand() command = %q, want %q", command, tt.want)
			}
			if len(args) == 0 || args[0] != tt.wantFirst {
				t.Fatalf("browserCommand() args = %#v, want first argument %q", args, tt.wantFirst)
			}
		})
	}
}

func TestBrowserCommandUnsupportedLinuxWhenNoOpener(t *testing.T) {
	// An unknown OS follows the Unix branch. In normal development machines
	// one of the openers exists; this assertion only verifies the error remains
	// actionable if a minimal Linux image has none of them.
	_, _, err := browserCommand("linux", "http://127.0.0.1:8080")
	if err != nil && !strings.Contains(err.Error(), "no supported browser opener") {
		t.Fatalf("unexpected error: %v", err)
	}
}
