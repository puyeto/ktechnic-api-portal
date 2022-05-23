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
	MeterName        string              `gorm:"not null" json:"meter_name"`
	MeterSerial      string              `gorm:"not null" json:"meter_serial"`
	Status           int8                `gorm:"not null;" json:"status" db:"status"`
	MeterDescription string              `gorm:"null" json:"meter_description"`
	CreatedAt        time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	AddedBy          uint32              `gorm:"-" json:"added_by"`
	Gateway          Gateway             `gorm:"-" json:"gateway_details"`
	Company          CompanyShortDetails `gorm:"-" json:"company_details"`
}

// Prepare ...
func (p *Meter) Prepare() {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.Status = 1
}

// Validate ...
func (p *Meter) Validate() error {

	if p.MeterName == "" {
		return errors.New("Meter Name is Required")
	}
	if p.MeterSerial == "" {
		return errors.New("Meter Serial is Required")
	}
	if p.GatewayID == 0 {
		return errors.New("Gateway is Required")
	}

	return nil
}

// SaveMeter ...
func (p *Meter) SaveMeter(db *gorm.DB) (*Meter, error) {
	var err error
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return &Meter{}, err
	}

	if err = tx.Debug().Model(&Meter{}).Create(&p).Error; err != nil {
		tx.Rollback()
		return &Meter{}, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &Meter{}, err
	}

	return p, nil
}

// ListAllMeters ...
func (p *Meter) ListAllMeters(db *gorm.DB) (*[]Meter, error) {
	var err error
	meters := []Meter{}
	var a CompanyShortDetails
	p.Gateway.Company = a
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return &meters, err
	}

	if p.CompanyID > 0 {
		err = tx.Debug().Where("company_id = ?", p.CompanyID).Model(&Meter{}).Limit(100).Find(&meters).Error
	} else {
		err = tx.Debug().Model(&Meter{}).Limit(100).Find(&meters).Error
	}

	if err != nil {
		tx.Rollback()
		return &meters, err
	}

	if len(meters) > 0 {
		for i := range meters {
			tx.Debug().Model(&Companies{}).Where("id = ?", meters[i].CompanyID).Take(&meters[i].Company)
			tx.Debug().Model(&Gateway{}).Where("id = ?", meters[i].GatewayID).Take(&meters[i].Gateway)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &meters, err
	}

	return &meters, nil
}

// FindMeterByID ...
func (p *Meter) FindMeterByID(db *gorm.DB, pid uint64) (*Meter, error) {
	var err error
	err = db.Debug().Model(&Meter{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return p, err
	}
	if p.ID != 0 {
		if err := db.Debug().Model(&User{}).Where("id = ?", p.CompanyID).Take(&p.Company).Error; err != nil {
			return p, err
		}
	}
	return p, nil
}

// UpdateAMeter ...
func (p *Meter) UpdateAMeter(db *gorm.DB) (*Meter, error) {

	var err error
	db.Debug().Model(&Meter{}).Where("id = ?", p.ID).Take(&Meter{}).UpdateColumns(
		map[string]interface{}{
			"gateway_id":        p.GatewayID,
			"meter_name":        p.MeterName,
			"meter_serial":      p.MeterSerial,
			"meter_description": p.MeterDescription,
			"updated_at":        p.UpdatedAt,
		},
	)
	err = db.Debug().Model(&Meter{}).Where("id = ?", p.ID).Take(&p).Error
	if err != nil {
		return &Meter{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CompanyID).Take(&p.Company).Error
		if err != nil {
			return &Meter{}, err
		}
	}
	return p, nil
}

// DeleteAMeter ...
func (p *Meter) DeleteAMeter(db *gorm.DB, vid uint32) (int64, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return 0, err
	}

	err := tx.Debug().Model(&Meter{}).Where("id = ?", vid).Take(&p).Error
	if err != nil {
		return 0, err
	}

	if err = tx.Debug().Model(&Meter{}).Where("id = ?", vid).Delete(&Meter{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Meter not found")
		}
		return 0, db.Error
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return db.RowsAffected, nil
}

// CountMeterTelemetryByID ...
func (p *Meter) CountMeterTelemetryByID(db *mongo.Database, mid uint64, filterfrom, filterto uint64) int {
	filter := bson.D{}
	if filterfrom > 0 && filterto > 0 {
		filter = bson.D{{"datetimestamp", bson.D{{"$gte", filterfrom}}}, {"datetimestamp", bson.D{{"$lte", filterto}}}}
	}

	deviceid := strconv.Itoa(int(mid))
	count, err := Count(deviceid, filter, nil)
	fmt.Printf("count %v with error %v", count, err)
	return count
}

// FindMeterTelemetryByID ...
func (p *Meter) FindMeterTelemetryByID(db *mongo.Database, mid uint64, order string, offset, limit int, filterfrom, filterto uint64) (*[]DeviceData, error) {
	var Telemetry []DeviceData

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
		return &Telemetry, err
	}
	defer cur.Close(ctx)

	for cur.Next(context.Background()) {
		item := DeviceData{}
		err := cur.Decode(&item)
		if err != nil {
			continue
		}
		Telemetry = append(Telemetry, item)

		// fmt.Println("Found a document: ", item)

	}
	if err := cur.Err(); err != nil {
		return &Telemetry, err
	}

	return &Telemetry, err
}

// Count returns the number of trip records in the database.
func Count(deviceid string, filter primitive.D, opts *options.FindOptions) (int, error) {
	app.CreateIndexMongo("data_" + deviceid)
	collection := app.MongoDB.Collection("data_" + deviceid)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	count, err := collection.CountDocuments(ctx, filter, nil)
	return int(count), err
}
