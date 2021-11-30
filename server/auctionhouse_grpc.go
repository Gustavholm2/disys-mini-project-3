package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Gustavholm2/disys-mini-project-3/shared"
	"google.golang.org/grpc"
)

type AuctionhouseServer struct {
	shared.UnimplementedAuctionhouseServer
	// Data goes here
}

var (
	highestBid shared.BidAmount
	endTime time.Time
)

func newAuctionhouseServer() *AuctionhouseServer {
	return &AuctionhouseServer{}
}

func (s *AuctionhouseServer) Bid(ctx context.Context, bidAmount *shared.BidAmount) (*shared.Empty, error) {
	if (time.Now().Before(endTime)) {
		if (bidAmount.Amount > highestBid.Amount){
			highestBid = *bidAmount
		}
		fmt.Println(bidAmount.Amount)
		return &shared.Empty{}, nil
	}
	return &shared.Empty{}, nil
}
func (s *AuctionhouseServer) Result(ctx context.Context, empty *shared.Empty) (*shared.Outcome, error) {
	if (time.Now().Before(endTime)){
		logI(fmt.Sprintf("the highest bid is: %d, and the auction is still running", &highestBid.Amount));
		return &shared.Outcome{Bid: &highestBid, IsOver: false}, nil
	}
	logI(fmt.Sprintf("the highest bid is: %d, and the auction is over", &highestBid.Amount));
	return &shared.Outcome{Bid: &highestBid, IsOver: true}, nil
}

func startServer(address string) {
	logI(fmt.Sprintf("Starting server with address: %s", address))
	lis, err := net.Listen("tcp", address)
	if err != nil {
	logF(err.Error())
	endTime = time.Now().Add(4*time.Minute)
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
