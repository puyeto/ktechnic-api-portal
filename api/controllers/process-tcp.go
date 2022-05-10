package controllers

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
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
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}

		fmt.Print(temp)

		result := "Received\n"
		c.Write([]byte(string(result)))
	}
	c.Close()
}

// HandleRequest Handles incoming requests.
func HandleRequest(conn net.Conn) {

	byteData := make([]byte, 83)

	for {
		reqLen, err := conn.Read(byteData)
		if err != nil {
			if err != io.EOF {
				fmt.Println("End of file error:", err)
			}
			fmt.Println("Error reading:", err.Error(), reqLen)
			return
		}

		// return Response
		result := "Received byte size = " + strconv.Itoa(reqLen) + "\n"
		fmt.Println(result)
		conn.Write([]byte(string(result)))

		// if reqLen == 0 {
		// 	return // connection already closed by client
		// }

		// if reqLen > 0 {
		// 	byteRead := bytes.NewReader(byteData)

		// 	mb := make([]byte, reqLen)
		// 	n1, _ := byteRead.Read(mb)

		// 	processRequest(conn, mb, n1)

		// }
		// opsRate.Mark(1)
	}
}

func readNextBytes(conn net.Conn, number int) (int, []byte) {
	bytes := make([]byte, number)

	reqLen, err := conn.Read(bytes)
	if err != nil {
		if err != io.EOF {
			fmt.Println("End of file error:", err)
		}
		fmt.Println("Error reading:", err.Error(), reqLen)
	}

	return reqLen, bytes
}

func processRequest(conn net.Conn, b []byte, byteLen int) {
	clientJobs := make(chan models.DeviceData)
	go generateResponses(clientJobs)

	var deviceData models.DeviceData

	if byteLen < 82 {
		fmt.Println("Invalid Byte Length = ", byteLen)
		return
	}

	byteReader := bytes.NewReader(b)

	scode := processSeeked(byteReader, 5, 0)
	deviceData.SystemCode = string(scode)
	if deviceData.SystemCode != "LEMON" {
		return
	}

	mid := processSeeked(byteReader, 1, 5)
	deviceData.MessageID = int(mid[0])

	bc := processSeeked(byteReader, 1, 6)
	deviceData.ByteCount = int(bc[0])

	gid := processSeeked(byteReader, 4, 7)
	deviceData.GatewayID = binary.LittleEndian.Uint32(gid)

	did := processSeeked(byteReader, 4, 11)
	deviceData.DeviceID = binary.LittleEndian.Uint32(did)
	deviceData.MeterID = deviceData.DeviceID

	en1 := processSeeked(byteReader, 4, 15)
	deviceData.VoltageLine1 = float32frombytes(en1)

	en2 := processSeeked(byteReader, 4, 19)
	deviceData.VoltageLine2 = float32frombytes(en2)

	en3 := processSeeked(byteReader, 4, 23)
	deviceData.VoltageLine3 = float32frombytes(en3)

	edl := processSeeked(byteReader, 4, 27)
	deviceData.TotalLineCurrent = float32frombytes(edl)

	eds := processSeeked(byteReader, 4, 31)
	deviceData.SystemPower = float32frombytes(eds)

	edsp := processSeeked(byteReader, 4, 35)
	deviceData.SystemPowerFactor = float32frombytes(edsp)

	edsf := processSeeked(byteReader, 4, 39)
	deviceData.SystemFrequency = float32frombytes(edsf)

	res := processSeeked(byteReader, 4, 43)
	deviceData.ImportActiveEnergy = float32frombytes(res)

	res = processSeeked(byteReader, 4, 47)
	deviceData.ImportReactiveEnergy = float32frombytes(res)

	res = processSeeked(byteReader, 1, 75)
	deviceData.UTCTimeSeconds = int(res[0])

	res = processSeeked(byteReader, 1, 76)
	deviceData.UTCTimeMinutes = int(res[0])

	res = processSeeked(byteReader, 1, 77)
	deviceData.UTCTimeHours = int(res[0])

	res = processSeeked(byteReader, 1, 78)
	deviceData.UTCTimeDay = int(res[0])

	mn := processSeeked(byteReader, 1, 79)
	deviceData.UTCTimeMonth = int(mn[0])

	yr := processSeeked(byteReader, 2, 80)
	deviceData.UTCTimeYear = int(binary.LittleEndian.Uint16(yr))

	checksum := processSeeked(byteReader, 1, 82)
	deviceData.Checksum = int(checksum[0])

	deviceData.DateTime = time.Date(deviceData.UTCTimeYear, time.Month(deviceData.UTCTimeMonth), deviceData.UTCTimeDay, deviceData.UTCTimeHours, deviceData.UTCTimeMinutes, deviceData.UTCTimeSeconds, 0, time.UTC)
	deviceData.DateTimeStamp = deviceData.DateTime.Unix()
	deviceData.CreatedOn = time.Now()

	fmt.Println(deviceData)
	chks := make([]byte, 1)
	for i := 5; i < 81; i++ {
		chks[0] += b[i]
	}

	if chks[0] != checksum[0] {
		return
	}

	clientJobs <- deviceData

}

func processSeeked(byteReader *bytes.Reader, bytesize, seek int64) []byte {
	byteReader.Seek(seek, 0)
	val := make([]byte, bytesize)
	byteReader.Read(val)
	return val
}

func float32frombytes(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func generateResponses(deviceData chan models.DeviceData) {
	for {
		// Wait for the next job to come off the queue.
		d := <-deviceData

		// LogToRedis(clientJob.DeviceData)
		server.logToMongoDB(d)

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

func readInt32(data []byte) (ret int32) {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &ret)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	return ret
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
