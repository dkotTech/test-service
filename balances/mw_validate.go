package balances

import (
	"context"
	"test-service/helpers/validate"
)

type mwValidate struct {
	inner Service
}

func NewMWValidate(inner Service) Service {
	return &mwValidate{inner: inner}
}

func (m *mwValidate) CurrentOne(ctx context.Context, r GetCurrentOneRequest) (GetCurrentOneResponse, error) {
	if err := validate.MustValidate(ctx, &r); err != nil {
		return GetCurrentOneResponse{}, err
	}

	return m.inner.CurrentOne(ctx, r)
}

func (m *mwValidate) Deposit(ctx context.Context, r DepositRequest) (DepositResponse, error) {
	if err := validate.MustValidate(ctx, &r); err != nil {
		return DepositResponse{}, err
	}

	return m.inner.Deposit(ctx, r)
}

func (m *mwValidate) Withdraw(ctx context.Context, r WithdrawRequest) (WithdrawResponse, error) {
	if err := validate.MustValidate(ctx, &r); err != nil {
		return WithdrawResponse{}, err
	}

	return m.inner.Withdraw(ctx, r)
}
