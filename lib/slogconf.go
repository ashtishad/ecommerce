package lib

import (
	"log/slog"
	"path/filepath"
)

// getSlogConf constructs and returns a pointer to a slog.HandlerOptions struct.
// This function customizes the log configuration by setting the logging level to debug
// and stripping the full directory path from the source's filename. The customized options
// are encapsulated in a slog.HandlerOptions and returned as a pointer.
//
// Returns:
//   - *slog.HandlerOptions: Pointer to a slog.HandlerOptions struct containing the logging configurations.
func getSlogConf() *slog.HandlerOptions {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			sourceVal, ok := a.Value.Any().(*slog.Source)
			if !ok {
				return a
			}
			sourceVal.File = filepath.Base(sourceVal.File)
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
