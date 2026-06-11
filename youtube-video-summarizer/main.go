package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	outputPath := flag.String("o", "", "output file (default: stdout)")
	youtubeURL := flag.String("yt", "", "YouTube URL: download normal subtitles first, fallback to auto-generated")
	lang := flag.String("lang", "en", "subtitle language code (used with -yt)")
	flag.Parse()

	if *youtubeURL != "" {
		runYoutubeMode(*youtubeURL, *lang, *outputPath)
		return
	}

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [-o <output>] <input.vtt>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -yt <URL> [-lang <code>] [-o <output>]\n", os.Args[0])
		os.Exit(1)
	}

	runLocalMode(flag.Arg(0), *outputPath)
}

func runLocalMode(inputPath, outputPath string) {
	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	result, err := cleanVTT(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error cleaning VTT: %v\n", err)
		os.Exit(1)
	}
	writeOutput(result, outputPath)
}

// extractVideoID pulls the canonical 11-char video ID from common YouTube URL forms.
func extractVideoID(url string) (string, error) {
	// Covers: watch?v=, /shorts/, /embed/, youtu.be/
	re := regexp.MustCompile(`(?:youtube\.com/(?:watch\?v=|shorts/|embed/)|youtu\.be/)([a-zA-Z0-9_-]{11})`)
	m := re.FindStringSubmatch(url)
	if m == nil {
		return "", fmt.Errorf("cannot extract video ID from URL")
	}
	return m[1], nil
}

// runSubDownload invokes yt-dlp for subtitle download.
// if autoSub is true, uses --write-auto-sub; otherwise --write-sub (normal).
// Returns the yt-dlp exit error (nil on success).
func runSubDownload(url, lang, vttDir string, autoSub bool) error {
	writeFlag := "--write-sub"
	if autoSub {
		writeFlag = "--write-auto-sub"
	}

	args := []string{
		writeFlag,
		"--sub-lang", lang,
		"--skip-download",
		"-o", "%(id)s",
		"-P", vttDir,
		url,
	}
	cmd := exec.Command("yt-dlp", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runYoutubeMode(url, lang, outputPath string) {
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: yt-dlp not found in PATH. Install it first.\n")
		os.Exit(1)
	}

	videoID, err := extractVideoID(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot extract video ID from URL: %v\n", err)
		os.Exit(1)
	}

	vttDir := "_vtt"
	if err := os.MkdirAll(vttDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating _vtt directory: %v\n", err)
		os.Exit(1)
	}

	// Predicted subtitle file path (yt-dlp -o "%(id)s" => <videoID>.<lang>.vtt)
	vttFile := filepath.Join(vttDir, videoID+"."+lang+".vtt")

	// Phase 1: try normal (human-authored) subtitles first
	normalErr := runSubDownload(url, lang, vttDir, false)
	if normalErr != nil {
		fmt.Fprintf(os.Stderr, "Error running yt-dlp (normal subs): %v\n", normalErr)
		os.Exit(1)
	}
	if fi, statErr := os.Stat(vttFile); statErr == nil && fi.Size() > 0 {
		// Normal subs succeeded — proceed to output
	} else {
		// Phase 2: fallback to auto-generated subtitles
		autoErr := runSubDownload(url, lang, vttDir, true)
		if autoErr != nil {
			fmt.Fprintf(os.Stderr, "Error running yt-dlp (auto subs): %v\n", autoErr)
			os.Exit(1)
		}
		if fi, statErr := os.Stat(vttFile); statErr != nil || fi.Size() == 0 {
			fmt.Fprintf(os.Stderr, "Error: no subtitles found (tried normal and auto-generated)\n")
			os.Exit(1)
		}
	}

	absPath, _ := filepath.Abs(vttFile)

	if outputPath != "" {
		f, err := os.Open(vttFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening subtitle file: %v\n", err)
			os.Exit(1)
		}
		result, err := cleanVTT(f)
		f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error cleaning VTT: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(outputPath, []byte(result+"\n"), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println(absPath)
}

// cleanVTT reads VTT content from a reader, strips markup, deduplicates, and returns clean prose text.
func cleanVTT(r io.Reader) (string, error) {
	timestampRe := regexp.MustCompile(`-->`)
	inlineTagRe := regexp.MustCompile(`<[^>]+>`)

	scanner := bufio.NewScanner(r)
	var result []string
	lastKept := ""

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed == "WEBVTT" || strings.HasPrefix(trimmed, "Kind:") || strings.HasPrefix(trimmed, "Language:") {
			continue
		}
		if timestampRe.MatchString(line) {
			continue
		}
		cleaned := inlineTagRe.ReplaceAllString(line, "")
		cleaned = strings.TrimSpace(cleaned)
		if cleaned == "" {
			continue
		}
		if cleaned == lastKept {
			continue
		}
		result = append(result, cleaned)
		lastKept = cleaned
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(result, " "), nil
}

func writeOutput(data, outputPath string) {
	if outputPath != "" {
		if err := os.WriteFile(outputPath, []byte(data+"\n"), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(data)
	}
}
