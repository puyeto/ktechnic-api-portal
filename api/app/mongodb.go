package app

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB ...
var MongoDB *mongo.Database

// InitializeMongoDB Initialize MongoDB Connection
func InitializeMongoDB(dbURL, dbName string, logger *logrus.Logger) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect(ctx)

	logger.Printf("Mongo DB initialized %v", dbName)
	return client.Database(dbName)
}

// CreateIndexMongo create a mongodn index
func CreateIndexMongo(colName string) (string, error) {
	mod := mongo.IndexModel{
		Keys: bson.M{
			"datetimestamp": -1, // index in ascending order
		}, Options: nil,
	}
	collection := MongoDB.Collection(colName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return collection.Indexes().CreateOne(ctx, mod)
}

func LogToMongoDB(d DataPacket) error {
	collection := MongoDB.Collection("data_" + strconv.FormatInt(int64(d.MeterNumber), 10))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, d)
	CreateIndexMongo("data_" + strconv.FormatInt(int64(d.MeterNumber), 10))
	return err
}

// LoglastSeenMongoDB update last seen
func LoglastSeenMongoDB(d DataPacket) error {
	data := bson.M{
		"$set": bson.M{
			"last_seen_date": d.DateTime,
			"last_seen_unix": d.DateTimeStamp,
			"updated_at":     time.Now(),
		},
	}

	return upsert(data, d.MeterNumber, "a_device_lastseen")
}

func upsert(data bson.M, meterNo uint64, table string) error {
	collection := MongoDB.Collection(table)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": meterNo}, data, opts)

	return err
}

// #D#GW#Mn#Date;time;pwr;batt;gsm;tpu;mtr_reading;valve_status;firmware_version/r/n
// D = data identifier.
// GW = gateway number.
// Mn = meter number.
// Date = Date in UTC format, DDMMYY, if no data, NA is sent.
// Time = Time in UTC format, HHMMSS, if no data, NA is sent.
// Pwr = Logical value indicating if device is connected to power, 1 = connected, 0 = disconnected.
// Batt = battery charge 0- 100%.
// Gsm = GSM signal strength 0-31, should be converted to 0 – 100% scale.
// Tpu = numerical values showing data yet to be uploaded.
// Mtr_reading = numerical water meter reading (no decimal points).
// Valve status = Logical value indicating valve status, 1= open, 0 = closed.
// Firmware_version = numerical value showing firmware version.

type DataPacket struct {
	DataIdentifier  string    `json:"data_identifier"`
	GatewayNumber   uint64    `json:"gateway_number"`
	MeterNumber     uint64    `json:"meter_number"`
	Date            string    `json:"date"`
	Time            string    `json:"time"`
	PowerConnected  uint8     `json:"power_connected"`
	BatteryCharge   uint8     `json:"battery_charge"`
	GSMSignal       uint8     `json:"gsm_signal"`
	TPU             int       `json:"tpu"`
	MeterReading    uint32    `json:"meter_reading"`
	ValveStatus     uint8     `json:"valve_status"`
	FirmwareVersion uint64    `json:"firmware_version"`
	DateTime        time.Time `json:"date_time"`
	DateTimeStamp   int64     `json:"date_time_stamp,omitempty"`
	CreatedOn       time.Time `json:"created_on"`
}
