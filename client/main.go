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
	fmt.Print("HostAddress: ")
	var address string
	fmt.Scanln(&address)

	fmt.Print("Name (no spaces): ")
	var participantName string
	fmt.Scanln(&participantName)

	startClient(address, participantName)
	
}

func startClient(address string, name string) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		logF(err.Error())
	}

	defer conn.Close()
	
	client := shared.NewAuctionhouseClient(conn)

	go mainloop(client, name)

	fmt.Scanln()
}

func mainloop(client shared.AuctionhouseClient, name string) {
	var (
		inputCommand string
		inputValue int
	)
	for (true) {
		fmt.Scanln(&inputCommand, &inputValue)
		switch inputCommand[0] {
		case 'r':
			bid, err := client.Result(context.Background(), &shared.Empty{})
			if (err != nil){
				logE(err.Error())
			}
			if (bid.IsOver){
				logI(fmt.Sprintf("This auction is over, and the highest bid was %d by %s", bid.Bid.Amount, bid.Bid.Owner))
			} else {
				logI(fmt.Sprintf("This auction is still running, the highest bid is %d by %s", bid.Bid.Amount, bid.Bid.Owner))
			}
			
		case 'b':
			_, err := client.Bid(context.Background(), &shared.BidAmount{Amount: int32(inputValue), Owner: name})
			if (err != nil){
				logE(err.Error())
			}
			logI(fmt.Sprintf("Made bid of %d", inputValue))
		}
	}
}
