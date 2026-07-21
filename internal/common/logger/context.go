package logger

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

// WithLogger یک کپی از context به همراه نمونه logger داده شده برمی‌گرداند
func WithLogger(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

// FromContext نمونه logger موجود در context را استخراج می‌کند
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
