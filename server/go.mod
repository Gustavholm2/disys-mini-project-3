module github.com/Gustavholm2/disys-mini-project-3/server

go 1.17

replace github.com/Gustavholm2/disys-mini-project-3/shared => ../shared

require (
	github.com/Gustavholm2/disys-mini-project-3/shared v0.0.0-20211123115938-cb9411c5dca3
	google.golang.org/grpc v1.42.0
)

require (
	github.com/golang/protobuf v1.5.0 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
