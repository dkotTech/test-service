package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	serrors "test-service/errors"
)

// Decode helper decode json body from request to T
func Decode[T any](ctx context.Context, r *http.Request) (T, error) {
	var request T

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, fmt.Errorf("failed to decode: %w", err)
	}

	return request, nil
}

// Encode helper encode response T to json and send it to writer
// if response is nil, send a created status code
func Encode[T any](_ context.Context, w http.ResponseWriter, response *T) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if response == nil {
		w.WriteHeader(http.StatusCreated)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("could not encode response: %w", err)
	}

	return nil
}

// CallService helper wrap a do function and call it after Decode and Encode a response of do function
func CallService[Req any, Res any](w http.ResponseWriter, r *http.Request, do func(ctx context.Context, _ Req) (Res, error)) {
	ctx := r.Context()

	req, err := Decode[Req](ctx, r)
	if err != nil {
		SendError(ctx, w, err)
		return
	}

	res, err := do(ctx, req)
	if err != nil {
		SendError(ctx, w, err)
		return
	}

	_ = Encode[Res](ctx, w, &res)
}

// SendError helper send a service error in json format
func SendError(ctx context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	err = Encode(ctx, w, serrors.Unwrap(err))
	if err != nil {
		log.Println(err)
	}
}
