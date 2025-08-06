package balances

import (
	"context"
	"test-service/helpers"
)

type mwErrors struct {
	inner Service
}

func NewMWErrors(inner Service) Service {
	return &mwErrors{inner: inner}
}

func (m *mwErrors) CurrentOne(ctx context.Context, r GetCurrentOneRequest) (GetCurrentOneResponse, error) {
	return helpers.ToServiceErrorWrap(ctx, r, m.inner.CurrentOne)
}

func (m *mwErrors) Deposit(ctx context.Context, r DepositRequest) (DepositResponse, error) {
	return helpers.ToServiceErrorWrap(ctx, r, m.inner.Deposit)
}

func (m *mwErrors) Withdraw(ctx context.Context, r WithdrawRequest) (WithdrawResponse, error) {
	return helpers.ToServiceErrorWrap(ctx, r, m.inner.Withdraw)
}
