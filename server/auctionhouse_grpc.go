package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Gustavholm2/disys-mini-project-3/shared"
	"google.golang.org/grpc"
)

type AuctionhouseServer struct {
	shared.UnimplementedAuctionhouseServer
	// Data goes here
}

var (
	highestBid shared.BidAmount
)

func newAuctionhouseServer() *AuctionhouseServer {
	return &AuctionhouseServer{}
}

func (s *AuctionhouseServer) Bid(ctx context.Context, bidAmount *shared.BidAmount) (*shared.Empty, error) {
if (bidAmount.Amount > highestBid.Amount){
		highestBid = *bidAmount
	}
	fmt.Println(bidAmount.Amount)
	return &shared.Empty{}, nil
}
func (s *AuctionhouseServer) Result(ctx context.Context, empty *shared.Empty) (*shared.Outcome, error) {
return nil, status.Errorf(codes.Unimplemented, "method Result not implemented")
}

func startServer(address string) {
	logI(fmt.Sprint("Starting server withaddress: %s", address))
	lis, err := net.Listen("tcp", address)
	if err != nil {
	logF(err.Error())
	}

var opts []grpc.ServerOption
	// ServerOption goes here

	grpcServer := grpc.NewServer(opts...)

	shared.RegisterAuctionhouseServer(grpcServer, newAuctionhouseServer())

	logI("server starting...")
	err = grpcServer.Serve(lis)
	if err != nil {
	logF(err.Error())
	}
}
