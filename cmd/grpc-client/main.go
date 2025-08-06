package main

import (
	"context"
	"flag"
	"log"
	grpcbalances "test-service/balances/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	accountID = flag.String("account_id", "7f55f0f8-ebe2-4522-8a96-25a46509885f", "Account ID")
)

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := grpcbalances.NewBalancesServiceClient(conn)

	resp, err := client.CurrentOne(context.Background(),
		&grpcbalances.GetCurrentOneRequest{AccountId: *accountID},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Balance of %s: %.2f\n", resp.AccountId, resp.Balance)
}
