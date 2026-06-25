package logger

import (
	"log/slog"
	"os"
	"time"
)

func New() *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Value = slog.StringValue(time.Now().UTC().Format(time.RFC3339Nano))
				}
				return a
			},
		}),
	)
}

