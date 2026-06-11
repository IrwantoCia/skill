# Skill Collection

This repository contains a collection of **skills** and **CLIs** for [OpenCode](https://opencode.ai)-based agents as well as standalone use.

## Structure

Each folder at this level is one standalone project. A project can be:

- **Skill** — agent instructions (`SKILL.md` file)
- **CLI** — command-line tool (with `main.go`, `Makefile`, etc.)
- **Both** — a CLI that can also be installed as an agent skill

Every project has its own `README.md`, `project.md`, and `current.md`.

## Current Projects

| Project | Description |
|---------|-------------|
| [youtube-video-summarizer](/youtube-video-summarizer) | Go CLI for cleaning YouTube subtitle (.vtt) files into clean prose. Can be installed as an OpenCode skill. |

## License

MIT
