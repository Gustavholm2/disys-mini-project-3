# disys-mini-project-3

## How to start Servers
Navigate to _/server_

Run the following command in a terminal:

`go run .`

Enter the port you wish to use for the server in question when prompted

The server can be closed by pressing [Enter]


## How to start Clients
Navigate to _/client_

Run the following command in a terminal: 

`go run .`

Enter the ports of the servers you wish to use when prompted

### Client commands
The following commands can be given to the client:

- `bid <number>` (where `<number>` is the bid amount) To make a bid.
- `result` To see the current highes bid and whether the auction is over.
- `quit` To close the client and create a log.

You only need to write the first letter of the commands, so `b 10` would be the same as `bid 10`.
