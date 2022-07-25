package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

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

	defer l.Close()
	rand.Seed(time.Now().Unix())

	logger.Infof("Listening on %v:%v", CONNHOST, CONNPORT)

	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Handle connections in a new goroutine.
		go controllers.HandleConnection(conn)
	}
}
