package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/ktechnics/ktechnics-api/api"
	"github.com/ktechnics/ktechnics-api/api/app"
	"github.com/ktechnics/ktechnics-api/api/controllers"
	"github.com/sirupsen/logrus"
)

const (
	CONNHOST = "0.0.0.0"
	CONNPORT = 9001
	CONNTYPE = "tcp"
)

func main() {

	// create the logger
	logger := logrus.New()
	app.InitLogger(logger)
	app.MongoDB = app.InitializeMongoDB("mongodb://root:safcom2012@172.105.34.129:27017/?authSource=admin", "lectrotel_portal", logger)

	go api.Run(logger)

	// tcpConnection()
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(CONNPORT))
	if err != nil {
		return
	}

	// Listen for incoming connections.
	l, err := net.ListenTCP(CONNTYPE, tcpAddr)
	if err != nil {
		panic(err)
	}

	var connections []net.Conn
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

	logger.Infof("Listening on %v:%v", CONNHOST, CONNPORT)

	for {
		// Listen for an incoming connection.
		conn, err := l.AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				logger.Printf("accept temp err: %v", ne)
				// continue
			}

			logger.Printf("accept err: %v", err)
			return
		}
		log.Println("Client ", conn.RemoteAddr(), " connected")

		// Handle connections in a new goroutine.
		go controllers.HandleRequest(conn)

		connections = append(connections, conn)
		if len(connections)%1000 == 0 {
			fmt.Printf("total number of connections: %v", len(connections))
		}
	}
}
