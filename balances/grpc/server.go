package grpc

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"test-service/balances"
	"test-service/helpers"
)

type WalletServer struct {
	UnimplementedBalancesServiceServer
	svc balances.Service
	log *slog.Logger
}

func NewServer(log *slog.Logger, svc balances.Service) BalancesServiceServer {
	return &WalletServer{svc: svc, log: log}
}

func (w *WalletServer) CurrentOne(
	ctx context.Context,
	req *GetCurrentOneRequest,
) (*GetCurrentOneResponse, error) {
	helpers.LogInfoGRPC(ctx, w.log, "CurrentOne", "balances")

	id, err := uuid.Parse(req.GetAccountId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account_id")
	}

	resp, err := w.svc.CurrentOne(ctx, balances.GetCurrentOneRequest{AccountID: id})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &GetCurrentOneResponse{
		AccountId: resp.AccountID.String(),
		Balance:   resp.Balance,
	}, nil
}
