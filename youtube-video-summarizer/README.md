# YouTube Video Summarizer

[![skills.sh](https://skills.sh/b/IrwantoCia/skill)](https://skills.sh/IrwantoCia/skill)

A Go-based CLI for downloading YouTube subtitles (.vtt) and cleaning them into readable prose. Designed as an agent skill for [OpenCode](https://opencode.ai) and compatible with the [Agent Skills](https://agentskills.io) standard.

## Prerequisites

- **Go** ≥ 1.26.1 — to build/install the CLI
- **yt-dlp** (latest) — required for `-yt` mode. [Install guide](https://github.com/yt-dlp/yt-dlp#installation)

> Zero external Go dependencies — stdlib only.

## CLI Install

```bash
# Recommended — install production binary
go install github.com/IrwantoCia/skill/youtube-video-summarizer@latest
```

The binary is installed as `youtube-video-summarizer` in `$GOPATH/bin`. Ensure it's on your `$PATH`.

### Local Build

```bash
git clone https://github.com/IrwantoCia/skill.git
cd skill/youtube-video-summarizer
make build   # produces ./yt-summarizer (local alias)
make install # installs to $GOPATH/bin
```

### Verify

```bash
youtube-video-summarizer -h
```

## Usage

### YouTube Mode

```bash
# English subtitles (default)
youtube-video-summarizer -yt "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

# Specific language
youtube-video-summarizer -yt "https://youtu.be/dQw4w9WgXcQ" -lang id

# Save cleaned transcript to file
youtube-video-summarizer -yt "https://www.youtube.com/watch?v=dQw4w9WgXcQ" -lang en -o summary.txt
```

### Local Mode (existing `.vtt` file)

```bash
youtube-video-summarizer subtitle.vtt
youtube-video-summarizer -o output.txt subtitle.vtt
```

## Output Contract

### YouTube Mode (exit code 0)

| Stream | Content |
|--------|---------|
| **stdout** | Absolute path to downloaded `.vtt` file (single line, e.g., `/home/user/project/_vtt/dQw4w9WgXcQ.en.vtt`) |
| **stderr** | yt-dlp progress only |
| **`-o <file>`** | Cleaned transcript written to file (if specified) |

Side effect: `_vtt/` directory created in CWD; `.vtt` files retained.

### Local Mode

| Stream | Content |
|--------|---------|
| **stdout** | Cleaned transcript (markup stripped, deduplicated) |
| **`-o <file>`** | Output written to file (stdout is empty) |

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-yt` | `""` | YouTube URL (enables YouTube mode) |
| `-lang` | `"en"` | ISO 639-1 language code (only with `-yt`) |
| `-o` | `""` | Output file path for cleaned transcript |

Exit codes: `0` = success, `1` = error (message on stderr).

## Subtitle Priority

1. **Normal (human-authored)** — via `yt-dlp --write-sub`. Preferred for better accuracy and punctuation.
2. **Auto-generated** — via `yt-dlp --write-auto-sub`. Fallback if normal subs unavailable.

If neither is found, exits with code 1 and an error message.

## Skill Installation (for Agents)

### Via `npx skills` (recommended)

```bash
npx skills add IrwantoCia/skill --skill youtube-video-summarizer -a opencode
```

Uses the [Agent Skills CLI](https://github.com/vercel-labs/skills) (22k+ stars), supporting OpenCode and 67+ other agents.

### Manual (OpenCode)

```bash
# Global (all projects)
mkdir -p ~/.config/opencode/skills/youtube-video-summarizer
cp SKILL.md ~/.config/opencode/skills/youtube-video-summarizer/

# Per-project
mkdir -p .opencode/skills/youtube-video-summarizer
cp SKILL.md .opencode/skills/youtube-video-summarizer/
```

**Requirements**: `youtube-video-summarizer` binary in `$PATH`, `yt-dlp` installed for YouTube mode.

## Limitations

- **Subtitles required**: Cannot work without subtitles (normal or auto-generated).
- **YouTube only**: Supports `youtube.com` and `youtu.be` URLs only.
- **No audio transcription**: Reads existing subtitles only — no speech-to-text.
- **No language auto-detect**: Use `-lang` explicitly.
- **yt-dlp dependency**: `-yt` mode requires `yt-dlp` on the system.

## License

MIT
