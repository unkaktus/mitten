package main

import (
	"io"
	"strings"

	"github.com/acarl005/stripansi"
)

// ShellFindReader looks for the typical shell idicator (#,$,>)
// to appear at the end of stdout, then signalling by closing
// Found channel.
type ShellFindReader struct {
	reader      io.Reader
	builder     *strings.Builder
	teeReader   io.Reader
	Found       chan struct{}
	SkippedEcho chan struct{}
	Enabled     bool
}

func NewShellFindReader(r io.Reader) *ShellFindReader {
	builder := &strings.Builder{}
	return &ShellFindReader{
		reader:      r,
		builder:     builder,
		teeReader:   io.TeeReader(r, builder),
		Found:       make(chan struct{}),
		SkippedEcho: make(chan struct{}),
		Enabled:     true,
	}
}

func (sfr *ShellFindReader) Read(p []byte) (int, error) {
	select {
	case <-sfr.Found:
		if !sfr.Enabled {
			discard := make([]byte, len(p))
			_, err := sfr.teeReader.Read(discard)
			accumulated := strings.TrimSpace(
				sfr.builder.String(),
			)
			accumulated = stripansi.Strip(accumulated)
			if accumulated != "" {
				if strings.Contains(accumulated, "#mitten_marker#") {
					sfr.Enabled = true
					sfr.builder.Reset()
					close(sfr.SkippedEcho)
				}
			}
			return 0, err
		}
		return sfr.reader.Read(p)
	default:
		n, err := sfr.teeReader.Read(p)
		accumulated := sfr.builder.String()
		accumulated = stripansi.Strip(accumulated)
		accumulated = strings.TrimSpace(accumulated)

		if accumulated != "" {
			switch rune(accumulated[len(accumulated)-1]) {
			case '$', '#', '>', '%':
				// Skip the two last characters are the same (banners)
				if accumulated[len(accumulated)-1] == accumulated[len(accumulated)-2] {
					break
				}
				close(sfr.Found)
				sfr.Enabled = false
				sfr.builder.Reset()
			}
		}
		return n, err
	}
}
