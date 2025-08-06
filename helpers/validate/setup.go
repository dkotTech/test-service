package validate

import (
	"context"
	"github.com/dkotTech/startfast"
	"github.com/go-playground/validator/v10"
)

// Validator setup a lazy validator
var Validator = startfast.NewEager(func() (*validator.Validate, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate, nil
})

// MustValidate handle struct for validation, return treasury error
func MustValidate[T any](ctx context.Context, t *T) error {
	err := Validator.MustGet().Struct(t)
	if err != nil {
		return err
	}
	return nil
}
