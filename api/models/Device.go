package models

import (
	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type DeviceData struct {
	ID                   string    `json:"id" bson:"_id,omitempty"`
	SystemCode           string    `json:"system_code"`
	MessageID            int       `json:"message_id"`
	ByteCount            int       `json:"byte_count"`
	GatewayID            uint32    `json:"gateway_id"`
	DeviceID             uint32    `json:"device_id"`
	MeterID              uint32    `json:"meter_id"`
	ModbusID             int       `json:"modbus_id"`
	EnergyReading        float32   `json:"energy_reading"`
	VoltageLine1         float32   `json:"voltage_line_1"`
	VoltageLine2         float32   `json:"voltage_line_2"`
	VoltageLine3         float32   `json:"voltage_line_3"`
	TotalLineCurrent     float32   `json:"total_line_current"`
	SystemPower          float32   `json:"system_power"`
	SystemPowerFactor    float32   `json:"system_power_factor"`
	SystemFrequency      float32   `json:"system_frequency"`
	ImportActiveEnergy   float32   `json:"import_active_energy"`
	ImportReactiveEnergy float32   `json:"import_reactive_energy"`
	UTCTimeSeconds       int       `json:"utc_time_seconds,omitempty"` // 1 byte
	UTCTimeMinutes       int       `json:"utc_time_minutes,omitempty"` // 1 byte
	UTCTimeHours         int       `json:"utc_time_hours,omitempty"`   // 1 byte
	UTCTimeDay           int       `json:"utc_time_day,omitempty"`     // 1 byte
	UTCTimeMonth         int       `json:"utc_time_month,omitempty"`   // 1 byte
	UTCTimeYear          int       `json:"utc_time_year,omitempty"`    // 1 byte
	DateTime             time.Time `json:"date_time,omitempty"`
	DateTimeStamp        int64     `json:"date_time_stamp,omitempty"`
	Checksum             int       `json:"checksum,omitempty"`
	Online               int8      `json:"online"`
	RelayStatus          int8      `json:"relay_status"`
	CreatedOn            time.Time `json:"created_on"`
}

// // DeviceData ...
// type DeviceData struct {
// 	ID                             string `json:"id" bson:"_id,omitempty"`
// 	SystemCode                     string
// 	MessageID                      int
// 	ByteCount                      int       `json:"byte_count"`
// 	GatewayID                      uint32    `json:"gateway_id"`
// 	DeviceID                       uint32    `json:"device_id"`
// 	MeterID                        uint32    `json:"meter_id"`
// 	EnergyDataLine1                float32   `json:"energy_data_line1,omitempty"`
// 	EnergyDataLine2                float32   `json:"energy_data_line2,omitempty"`
// 	EnergyDataLine3                float32   `json:"energy_data_line3,omitempty"`
// 	EnergyDataTotalLineCurrent     float32   `json:"energy_data_total_line_current,omitempty"`
// 	EnergyDataSystemPower          float32   `json:"energy_data_system_power,omitempty"`
// 	EnergyDataSystemPowerFactor    float32   `json:"energy_data_system_power_factor,omitempty"`
// 	EnergyDataSystemPowerFrequency float32   `json:"energy_data_system_power_frequency,omitempty"`
// 	EnergyDataImportActiveEnergy   float32   `json:"energy_data_import_active_energy,omitempty"`
// 	EnergyDataImportReactiveEnergy float32   `json:"energy_data_import_reactive_energy,omitempty"`
// 	UTCTimeSeconds                 int       `json:"utc_time_seconds,omitempty"` // 1 byte
// 	UTCTimeMinutes                 int       `json:"utc_time_minutes,omitempty"` // 1 byte
// 	UTCTimeHours                   int       `json:"utc_time_hours,omitempty"`   // 1 byte
// 	UTCTimeDay                     int       `json:"utc_time_day,omitempty"`     // 1 byte
// 	UTCTimeMonth                   int       `json:"utc_time_month,omitempty"`   // 1 byte
// 	UTCTimeYear                    int       `json:"utc_time_year,omitempty"`    // 1 byte
// 	DateTime                       time.Time `json:"date_time,omitempty"`
// 	DateTimeStamp                  int64     `json:"date_time_stamp,omitempty"`
// 	Checksum                       int       `json:"checksum,omitempty"`
// }

// SaveDeviceData ...
func (d DeviceData) SaveDeviceData(db *mongo.Database) error {
	collection := db.Collection("data_" + strconv.FormatInt(int64(d.DeviceID), 10))
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.InsertOne(ctx, d)
	return err
}
