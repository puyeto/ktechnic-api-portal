package controllers

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ktechnics/ktechnics-api/api/app"
	"github.com/ktechnics/ktechnics-api/api/models"
)

const queueLimit = 50

var server = Server{}

func HandleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		dataPackers := make(chan app.DataPacket)
		go generateResponses(dataPackers)

		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}

		result := "Received\n"
		c.Write([]byte(string(result)))

		// Process Data
		processedData, err := processRequest(temp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		dataPackers <- processedData
	}

	c.Close()
}

func processRequest(dataPacket string) (app.DataPacket, error) {
	var data = app.DataPacket{}
	// #D#867857030306398#49#090522;140704;1;100;30;8;0;1;24
	dataPacketArray := splitString(trimLeftChar(dataPacket), ";")
	meterdetails := splitString(dataPacketArray[0], "#")

	data.DataIdentifier = meterdetails[0]
	if data.DataIdentifier != "D" {
		return data, errors.New("Invalid Data")
	}

	data.GatewayNumber = stringToUInt64(meterdetails[1])
	if data.GatewayNumber == 0 {
		return data, errors.New("Invalid Gateway")
	}
	data.MeterNumber = stringToUInt64(meterdetails[2])

	d := splitString(meterdetails[3], "")
	t := splitString(dataPacketArray[1], "")
	myDateString := "20" + d[4] + d[5] + "-" + d[2] + d[3] + "-" + d[0] + d[1] + " " + t[0] + t[1] + ":" + t[2] + t[3] + ":" + t[4] + t[5]
	myDate, err := time.Parse("2006-01-02 15:04:05", myDateString)
	if err != nil {
		return data, err
	}

	data.Date = myDate.Format("2006-01-02")
	data.Time = myDate.Format("15:04:05")
	data.PowerConnected = uint8(stringToUInt64(dataPacketArray[2]))
	data.BatteryCharge = uint8(stringToUInt64(dataPacketArray[3]))
	data.GSMSignal = uint8(stringToUInt64(dataPacketArray[4]))
	data.TPU = int(stringToUInt64(dataPacketArray[5]))
	data.MeterReading = uint32(stringToUInt64(dataPacketArray[6]))
	data.ValveStatus = uint8(stringToUInt64(dataPacketArray[7]))
	data.FirmwareVersion = stringToUInt64(dataPacketArray[8])
	data.DateTime = myDate
	data.DateTimeStamp = myDate.Unix()
	created := time.Now()
	data.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", created.Format("2006-01-02 15:04:05"))

	return data, nil
}

func stringToUInt64(val string) uint64 {
	n, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func splitString(str, sep string) []string {
	return strings.Split(str, sep)
}

func generateResponses(dataPackers chan app.DataPacket) {

	for {
		clientJob := <-dataPackers
		fmt.Println("processed : ", clientJob)

				// use a WaitGroup
				var wg sync.WaitGroup

				// Wait for the next job to come off the queue.
				// LogToRedis(clientJob)

				// make a channel with a capacity of 100.
				jobChan := make(chan app.DataPacket, queueLimit)

				worker := func(jobChan <-chan app.DataPacket) {
					defer wg.Done()
					for job := range jobChan {
						// updateVehicleStatus(job.DeviceID, "online", "Online", job.DateTime)

						// SaveAllData(job)
						if err := app.LogToMongoDB(job); err != nil {
							fmt.Printf("Mongo DB - logging error: %v", err)
						}
						if err := app.LoglastSeenMongoDB(job); err != nil {
							fmt.Printf("Mongo DB - logging last seen error: %v", err)
						}

						// meterid = processedData.MeterNumber
						// devicetime = processedData.DateTime

						// updateVehicleStatus(meterid, "online", "Online", devicetime)
					}
				}

				// increment the WaitGroup before starting the worker
				wg.Add(1)
				go worker(jobChan)

				// enqueue a job
				jobChan <- clientJob

				// to stop the worker, first close the job channel
				close(jobChan)

				// then wait using the WaitGroup
				WaitTimeout(&wg, 1*time.Second)
	}
}

// WaitTimeout does a Wait on a sync.WaitGroup object but with a specified
// timeout. Returns true if the wait completed without timing out, false
// otherwise.
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

// FloatToString ...
func FloatToString(inputnum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputnum, 'f', 6, 64)
}

// logToMongoDB ...
func (s *Server) logToMongoDB(deviceData models.DeviceData) error {
	return deviceData.SaveDeviceData(app.MongoDB)
}
