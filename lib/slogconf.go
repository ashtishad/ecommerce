package lib

import (
	"log/slog"
	"path/filepath"
)

// GetSlogConf constructs and returns a pointer to a slog.HandlerOptions struct.
// This function customizes the log configuration by setting the logging level to debug
// and stripping the full directory path from the source's filename. The customized options
// are encapsulated in a slog.HandlerOptions and returned as a pointer.
//
// Returns:
//   - *slog.HandlerOptions: Pointer to a slog.HandlerOptions struct containing the logging configurations.
func GetSlogConf() *slog.HandlerOptions {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}

	handlerOpts := slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: replace,
	}

	return &handlerOpts
}
