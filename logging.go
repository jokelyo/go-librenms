package librenms

import (
	"log/slog"
	"net/http"
)

// logRequestAttr creates a slog.Attr for logging HTTP request details.
func logRequestAttr(req *http.Request) slog.Attr {
	if req == nil {
		return slog.String("request", "nil")
	}
	return slog.Group("request",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
	)
}

// logResponseAttr creates a slog.Attr for logging HTTP response details.
func logResponseAttr(resp *http.Response) slog.Attr {
	if resp == nil {
		return slog.String("response", "nil")
	}
	return slog.Group("response",
		slog.Int("status", resp.StatusCode),
		slog.String("status_text", resp.Status),
		slog.String("content_type", resp.Header.Get("Content-Type")),
	)
}
