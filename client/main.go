package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/Gustavholm2/disys-mini-project-3/shared"
)

var (
	// Logs str
	logI func(str string)
	// Logs str as an error
	logE func(str string)
	// Logs str a fatal error, then panics
	logF func(str string)
)

func main() {
	var buf bytes.Buffer
	var loggerInf = log.New(&buf, "LOG|INFO: ", log.Lshortfile|log.Lmicroseconds)
	var loggerErr = log.New(&buf, "LOG|ERR: ", log.Lshortfile|log.Lmicroseconds)
	var loggerFat = log.New(&buf, "LOG|FATAL: ", log.Lshortfile|log.Lmicroseconds)
	defer func() {
		fmt.Println(&buf)
		os.WriteFile("log.txt", buf.Bytes(), 0644)
	}()
	logI = func(str string) {
		loggerInf.Output(2, str)
		fmt.Println(str)
	}
	logE = func(str string) {
		loggerErr.Output(2, str)
		fmt.Println(str)
	}
	logF = func(str string) {
		loggerFat.Output(2, str)
		panic(str)
	}

	logI(fmt.Sprintf("=== LOG START: %v ===", time.Now().Format(time.RFC1123)))
	fmt.Print("NumberOfConnections: ")
	var numberOfConnections int
	fmt.Scanln(&numberOfConnections)

	serverAddresses := make([]string, numberOfConnections)

	for i := 0; i < numberOfConnections; i++ {
		fmt.Print("HostAddress: ")
		fmt.Scanln(&serverAddresses[i])
	}

	fmt.Print("Name (no spaces): ")
	var participantName string
	fmt.Scanln(&participantName)

	startClient(serverAddresses, participantName)

}

func startClient(addresses []string, name string) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conns := make([]*grpc.ClientConn, len(addresses))

	for i, address := range addresses {
		conn, err := grpc.Dial(address, opts...)
		if err != nil {
			logF(err.Error())
		}
		conns[i] = conn
	}

	defer func() {
		for _, conn := range conns {
			conn.Close()
		}
	}()

	clients := make([]shared.AuctionhouseClient, len(conns))

	for i, conn := range conns {
		clients[i] = shared.NewAuctionhouseClient(conn)
	}

	mainloop(clients, name)
}

func mainloop(clients []shared.AuctionhouseClient, name string) {
	var (
		inputCommand string
		inputValue   int
		quit         = false
	)
	for !quit {
		fmt.Print(">")
		fmt.Scanln(&inputCommand, &inputValue)
		switch inputCommand[0] {
		case 'r':
			bids := make([]*shared.Outcome, len(clients))
			for i, c := range clients {
				bid, err := c.Result(context.Background(), &shared.Empty{})
				if err != nil {
					logE(err.Error())
				}
				bids[i] = bid
			}
			bid, err := GetResultConsensus(bids, len(clients))
			if err != nil {
				logE(err.Error())
			}

			if bid.IsOver {
				logI(fmt.Sprintf("This auction is over, and the highest bid was %d by %s", bid.Bid.Amount, bid.Bid.Owner))
			} else {
				logI(fmt.Sprintf("This auction is still running, the highest bid is %d by %s", bid.Bid.Amount, bid.Bid.Owner))
			}

		case 'b':
			logI(fmt.Sprintf("Making bid of %d", inputValue))
			for _, c := range clients {
				_, err := c.Bid(context.Background(), &shared.BidAmount{Amount: int32(inputValue), Owner: name})
				if err != nil {
					logE(err.Error())
				}
			}
		default:
			fmt.Printf("Quit? [y/N]> ")
			var confirmation string
			fmt.Scanln(&confirmation)
			if confirmation == "y" {
				quit = true
			}
		}
	}
}

func GetResultConsensus(outcomes []*shared.Outcome, totalNodes int) (*shared.Outcome, error) {
	votes := make(map[string]int)
	for _, outcome := range outcomes {
		str := fmt.Sprint(outcome)
		votes[str]++
		if float32(votes[str]) > float32(totalNodes)/2 {
			return outcome, nil
		}
	}
	return nil, fmt.Errorf("no consensus")
}
