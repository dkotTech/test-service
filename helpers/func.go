package helpers

import (
	"context"
	"log/slog"
	"net/http"
	serrors "test-service/errors"
)

// ToServiceErrorWrap helper to wrap a service for return a service error type
func ToServiceErrorWrap[Req any, Res any](ctx context.Context, req Req, f func(ctx context.Context, _ Req) (Res, error)) (Res, error) {
	res, err := f(ctx, req)
	if err != nil {
		return res, serrors.NewServiceError(ctx, err.Error(), serrors.UserVisible)
	}

	return res, err
}

// LogInfoHTTP helper log info from http request
func LogInfoHTTP(log *slog.Logger, r *http.Request, msg, service string) {
	log.InfoContext(r.Context(), msg, "service", service, "method", r.Method)
}

// LogInfoGRPC helper log info from gRPC request
func LogInfoGRPC(ctx context.Context, log *slog.Logger, msg, service string) {
	log.InfoContext(ctx, msg, "service", service)
}
