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

const (
	auctionTime = 30 // in seconds
)

var (
	highestBid shared.BidAmount
	endTime    time.Time
)

func newAuctionhouseServer() *AuctionhouseServer {
	return &AuctionhouseServer{}
}

func (s *AuctionhouseServer) Bid(ctx context.Context, bidAmount *shared.BidAmount) (*shared.Empty, error) {
	if highestBid.Amount == -1 {
		endTime = time.Now().Add(auctionTime * time.Second)
		logI(fmt.Sprintf("Auction started at with starting bid of %d. Closing at %v", highestBid.Amount, endTime.Format("15:04:05")))
	}
	if time.Now().Before(endTime) {
		logI(fmt.Sprintf("%s makes a bid of %d", bidAmount.Owner, bidAmount.Amount))
		if bidAmount.Amount > highestBid.Amount {
			highestBid = *bidAmount
			logI(fmt.Sprintf("New highest bid of %d by %s", highestBid.Amount, highestBid.Owner))
		} else {
			logI("Bid not higher than highest bid. Ignoring")
			return nil, fmt.Errorf("Bid too low (current highest bid is %d)", highestBid.Amount)
		}
		return &shared.Empty{}, nil
	} else {
		return nil, fmt.Errorf("the auction is over. no more bids can be made")
	}
}

func (s *AuctionhouseServer) Result(ctx context.Context, empty *shared.Empty) (*shared.Outcome, error) {
	if time.Now().Before(endTime) {
		logI(fmt.Sprintf("the highest bid is: %d, and the auction is still running", highestBid.Amount))
		return &shared.Outcome{Bid: &highestBid, IsOver: false}, nil
	}
	logI(fmt.Sprintf("the highest bid is: %d, and the auction is over", highestBid.Amount))
	return &shared.Outcome{Bid: &highestBid, IsOver: true}, nil
}

func startServer(address string) {
	logI(fmt.Sprintf("Starting server with address: %s", address))
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logF(err.Error())
	}

	var opts []grpc.ServerOption
	// ServerOption goes here

	grpcServer := grpc.NewServer(opts...)

	shared.RegisterAuctionhouseServer(grpcServer, newAuctionhouseServer())

	highestBid = shared.BidAmount{Amount: -1, Owner: "AUCTION NOT STARTED"}

	logI("server starting...")
	err = grpcServer.Serve(lis)
	if err != nil {
		logF(err.Error())
	}
}
