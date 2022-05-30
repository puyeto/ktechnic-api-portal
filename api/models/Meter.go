package models

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ktechnics/ktechnics-api/api/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Meter ...
type Meter struct {
	ID               uint32              `gorm:"primary_key;auto_increment" json:"id"`
	CompanyID        uint32              `gorm:"not null;" json:"company_id"`
	GatewayID        uint32              `gorm:"not null;" json:"gateway_id"`
	PricePlanID      uint32              `json:"price_plan_id" gorm:"not null;"`
	MeterName        string              `gorm:"not null;" json:"meter_name"`
	MeterNumber      uint64              `gorm:"not null;" json:"meter_number"`
	Status           int8                `gorm:"not null;" json:"status" db:"status"`
	MeterDescription string              `gorm:"null;" json:"meter_description"`
	ValveStatus      int8                `gorm:"not null;" json:"valve_status"`
	LastSeen         time.Time           `gorm:"null;" json:"last_seen"`
	CreatedAt        time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	AddedBy          uint32              `gorm:"not null;" json:"added_by"`
	Gateway          Gateway             `gorm:"-" json:"gateway_details"`
	Company          CompanyShortDetails `gorm:"-" json:"company_details"`
}

// Prepare ...
func (m *Meter) Prepare() {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.ValveStatus = 1
	m.Status = 1
}

// Validate ...
func (m *Meter) Validate() error {

	if m.MeterName == "" {
		return errors.New("Meter Name is Required")
	}
	if m.MeterNumber == 0 {
		return errors.New("Meter Number is Required")
	}
	if m.GatewayID == 0 {
		return errors.New("Gateway is Required")
	}
	if m.PricePlanID == 0 {
		return errors.New("Price Plan is Required")
	}

	return nil
}

// SaveMeter ...
func (m *Meter) SaveMeter(db *gorm.DB) (*Meter, error) {
	exist, err := m.IsMeterExist(db)
	if err != nil {
		return m, err
	}
	if exist == true {
		return m, errors.New("Meter already exist")
	}

	fmt.Println(m.AddedBy)

	if err := db.Debug().Model(m).Create(&m).Error; err != nil {
		return m, err
	}

	db.Debug().Table("companies").Model(&Companies{}).Where("id = ?", m.CompanyID).Take(&m.Company)
	return m, nil
}

func (m *Meter) IsMeterExist(db *gorm.DB) (bool, error) {
	var result struct {
		Found bool
	}
	err := db.Raw("SELECT EXISTS(SELECT 1 FROM meters WHERE meter_number = ?) AS found", m.MeterNumber).Scan(&result).Error
	return result.Found, err
}

// Count Meters ...
func (m *Meter) CountMeters(db *gorm.DB, roleid uint32) int {
	var count int
	query := db.Debug().Model(m)
	if roleid == 1002 {
		query = query.Where("company_id = ?", m.CompanyID)
	} else if roleid > 1002 {
		query = query.Where("added_by = ?", m.AddedBy)
	}

	query.Count(&count)
	return count
}

// ListAllMeters ...
func (m *Meter) ListAllMeters(db *gorm.DB, roleid uint32, offset, limit int) (*[]Meter, error) {
	var err error
	meters := []Meter{}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return &meters, err
	}

	query := tx.Debug()
	if roleid == 1002 {
		query = query.Where("company_id = ?", m.CompanyID)
	} else if roleid > 1002 {
		query = query.Where("added_by = ?", m.AddedBy)
	}

	err = query.Model(m).Offset(offset).Limit(limit).Find(&meters).Error

	if err != nil {
		return &meters, err
	}

	if len(meters) > 0 {
		for i := range meters {
			tx.Debug().Table("companies").Model(&Companies{}).Where("id = ?", meters[i].CompanyID).Take(&meters[i].Company)
			tx.Debug().Model(&Gateway{}).Where("id = ?", meters[i].GatewayID).Take(&meters[i].Gateway)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return &meters, err
	}

	return &meters, nil
}

// FindMeterByID ...
func (m *Meter) GetMeterByID(db *gorm.DB) (*Meter, error) {
	var err error
	err = db.Debug().Model(m).Where("id = ?", m.ID).Take(&m).Error
	if err != nil {
		return m, err
	}
	if m.ID != 0 {
		db.Debug().Table("companies").Model(&CompanyShortDetails{}).Where("id = ?", m.CompanyID).Take(&m.Company)
		db.Debug().Model(&Gateway{}).Where("id = ?", m.GatewayID).Take(&m.Gateway)
	}
	return m, nil
}

// GetMeterByMeterNumber ...
func (m *Meter) GetMeterByMeterNumber(db *gorm.DB) (*Meter, error) {
	if err := db.Debug().Model(m).Where("meter_number = ?", m.MeterNumber).Take(&m).Error; err != nil {
		return m, err
	}

	if m.ID != 0 {
		db.Debug().Table("companies").Model(&CompanyShortDetails{}).Where("id = ?", m.CompanyID).Take(&m.Company)
		db.Debug().Model(&Gateway{}).Where("id = ?", m.GatewayID).Take(&m.Gateway)
	}
	return m, nil
}

// UpdateAMeter ...
func (m *Meter) UpdateAMeter(db *gorm.DB) (*Meter, error) {
	err := db.Debug().Model(m).Where("id = ?", m.ID).Updates(
		map[string]interface{}{
			"gateway_id":        m.GatewayID,
			"price_plan_id":     m.PricePlanID,
			"meter_name":        m.MeterName,
			"meter_number":      m.MeterNumber,
			"meter_description": m.MeterDescription,
			"updated_at":        m.UpdatedAt,
			"status":            m.Status,
			"valve_status":      m.ValveStatus,
			"last_seen":         m.LastSeen,
		},
	).Error

	return m, err
}

// DeleteAMeter ...
func (m *Meter) DeleteAMeter(db *gorm.DB, vid uint32) error {

	if err := db.Debug().Model(m).Where("id = ?", vid).Delete(m).Error; err != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return errors.New("Meter not found")
		}
		return db.Error
	}

	return nil
}

// CountMeterTelemetryByID ...
func (m *Meter) CountMeterTelemetryByID(db *mongo.Database, mid uint64, filterfrom, filterto uint64) int {
	filter := bson.D{}
	if filterfrom > 0 && filterto > 0 {
		filter = bson.D{{"datetimestamp", bson.D{{"$gte", filterfrom}}}, {"datetimestamp", bson.D{{"$lte", filterto}}}}
	}

	deviceid := strconv.Itoa(int(mid))
	count, _ := Count(deviceid, filter, nil)
	// fmt.Printf("count %v with error %v", count, err)
	return count
}

// FindMeterTelemetryByID ...
func (m *Meter) FindMeterTelemetryByID(db *mongo.Database, mid uint64, order string, offset, limit int, filterfrom, filterto uint64) (*[]app.DataPacket, error) {
	var telemetry []app.DataPacket

	// Get collection
	collection := db.Collection("data_" + strconv.FormatInt(int64(mid), 10))
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	app.CreateIndexMongo("data_" + strconv.FormatInt(int64(mid), 10))

	findOptions := options.Find()
	// Sort by `price` field descending
	if order == "asc" {
		findOptions.SetSort(bson.D{{"datetimestamp", 1}})
	} else {
		findOptions.SetSort(bson.D{{"datetimestamp", -1}})
	}

	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	filter := bson.D{}
	if filterfrom > 0 && filterto > 0 {
		filter = bson.D{{"datetimestamp", bson.D{{"$gte", filterfrom}}}, {"datetimestamp", bson.D{{"$lte", filterto}}}}
	}

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return &telemetry, err
	}
	defer cur.Close(ctx)

	for cur.Next(context.Background()) {
		item := app.DataPacket{}
		err := cur.Decode(&item)
		if err != nil {
			continue
		}
		telemetry = append(telemetry, item)

		// fmt.Println("Found a document: ", item)

	}
	if err := cur.Err(); err != nil {
		return &telemetry, err
	}

	return &telemetry, err
}

// Count returns the number of trip records in the database.
func Count(deviceid string, filter primitive.D, opts *options.FindOptions) (int, error) {
	app.CreateIndexMongo("data_" + deviceid)
	collection := app.MongoDB.Collection("data_" + deviceid)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	count, err := collection.CountDocuments(ctx, filter, nil)
	return int(count), err
}
