BINARY   := yt-summarizer
MODULE   := github.com/IrwantoCia/skill/youtube-video-summarizer

.PHONY: all build install run clean help

all: build

build:
	go build -o $(BINARY) .

install:
	go install .

run: build
	./$(BINARY) $(ARGS)

clean:
	rm -f $(BINARY)

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  all              Build binary (default)"
	@echo "  build     Build local binary as ./$(BINARY)"
	@echo "  install   Install to \$$(GOPATH)/bin as youtube-video-summarizer"
	@echo "  run       Build + run with ARGS='...'"
	@echo "  clean     Remove $(BINARY)"
	@echo "  help      Show this message"
