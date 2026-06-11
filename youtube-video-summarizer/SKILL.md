---
name: youtube-video-summarizer
description: Summarize YouTube video content by fetching and cleaning transcript text via CLI, then answering user questions.
compatibility: opencode >= 1.0
allowed-tools:
  - bash
  - read
  - write
  - edit
  - glob
  - grep
---

# YouTube Video Summarizer — Agent Runtime Guide

## When to Use

Use this skill when the user's message contains **both**:

- A **question** about the content of a YouTube video.
- A **YouTube URL** (`youtube.com/watch?v=...`, `youtu.be/...`, etc.).

Do **not** use this skill for:

- General web search or non-YouTube URLs.
- Requests that only ask to "clean a transcript" without a content question.

## Input Extraction

| Field | Source | Required |
|-------|--------|----------|
| **YouTube URL** | First `youtube.com` or `youtu.be` link in message. Prefer `watch?v=` form. | Yes |
| **Question** | Remainder of message — the query to answer. | Yes |
| **Language hint** | Optional. If user mentions a language, map to ISO 639-1 code. | No |

URL normalization: strip tracking/playlist params; only `v=` is needed.

## Execution

1. Extract URL, question, and optional language code from user message.
2. Default language is `en`. If user explicitly mentions a language, use ISO 639-1 code.
3. Run the CLI:
   ```bash
   youtube-video-summarizer -yt "<normalized_url>" -lang <code>
   ```
   (The `youtube-video-summarizer` binary must be installed and in `$PATH`.)
4. If exit code 0: **stdout** contains the absolute path to a `.vtt` file (single line, no other text). Read that file to obtain the raw WebVTT transcript. Optionally use `-o <file>` to also write a cleaned (markup-stripped, deduplicated) prose version.
5. If exit code non-zero: read stderr and handle per [Fallback](#fallback-behavior).
6. Answer the user's question grounded in the transcript text.

### Language Handling

- Default language code is `en`.
- If user mentions a language, pass the ISO 639-1 code (e.g., `id`, `es`).
- The CLI does **not** auto-detect. Infer from user's mention or video context.
- If transcript is garbled or wrong, retry once with a different likely code.

### Subtitle Preference (Normal → Auto Fallback)

The CLI tries subtitles in this priority order:
1. **Normal (human-authored) subtitles** — via `yt-dlp --write-sub`. Preferred when available because they have better accuracy and punctuation.
2. **Auto-generated subtitles** — via `yt-dlp --write-auto-sub`. Used as fallback only when normal subtitles are absent.

The agent does not need to handle this logic — it is transparent to the caller. On success, stdout always contains the absolute `.vtt` path regardless of which subtitle type was used.

### CLI Flags

| Flag | Default | Purpose |
|------|---------|---------|
| `-o` | `""` | Output file path for cleaned transcript (optional) |
| `-yt` | `""` | YouTube URL |
| `-lang` | `"en"` | Language code |

Exit codes: `0` = success (output on stdout), `1` = error (message on stderr).

### Output Contract (YouTube mode, exit 0)

| Stream | Content |
|--------|---------|
| stdout | Absolute path to downloaded `.vtt` file (e.g., `/home/user/project/_vtt/dQw4w9WgXcQ.en.vtt`) — one line only |
| stderr | Nothing (except yt-dlp progress) |
| `-o <file>` | If specified, cleaned transcript written there |

Side effect: `_vtt/` directory is created in CWD; the `.vtt` file is persisted and **not** cleaned up.

## Fallback Behavior

| Outcome | Agent Action |
|---------|--------------|
| Exit 0, path on stdout | Read `.vtt` file and answer. |
| Exit 0, file empty or missing | Report "Transcript retrieved but empty." Suggest different language or check video availability. |
| Exit 1, "yt-dlp not found" | Inform user: `yt-dlp` is not installed on this system. |
| Exit 1, "no subtitles found" | No subtitles available (normal or auto). Try a different language code or note the video may lack subtitles entirely. |
| Exit 1, "cannot extract video ID" | The URL is malformed or not a recognized YouTube URL. |
| Exit 1, other error | Report the stderr message concisely. Suggest a different video or using a local transcript. |
| Network / timeout | Inform user of transient failure; suggest retry. |

If transcript retrieval fails entirely, answer from general knowledge if appropriate, but clearly state you could not access the video transcript.

## Answer Style

- Ground your answer in the transcript text; reference specific parts.
- If the transcript is long, summarize first, then answer specifics.
- Use the same language as the user's question (or transcript, if more natural).
- Do not fabricate details not present in the transcript. If the transcript does not cover the asked topic, say so.

## Local `.vtt` Mode (Secondary)

For offline testing with local `.vtt` files:
```bash
youtube-video-summarizer transcript.vtt
youtube-video-summarizer -o out.txt transcript.vtt
```
Not part of the primary Q&A flow — use only for manual verification or debugging.
