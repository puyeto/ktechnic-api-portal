package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Odometers ...
type Odometers struct {
	ID            uint32    `gorm:"primary_key;auto_increment" json:"id"`
	VehicleID     uint32    `gorm:"not null" json:"vehicle_id"`
	OdometerValue float64   `gorm:"not null" json:"odometer_value,omitempty"`
	Units         string    `json:"units,omitempty"`
	DateAdded     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date_added"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	AddedBy       uint32    `gorm:"not null" json:"added_by"`
}

// InitializeOdometer ...
func InitializeOdometer(vehicleid uint32, odometervalue float64, dateadded time.Time, addedby uint32) *Odometers {
	var od Odometers
	od.VehicleID = vehicleid
	od.OdometerValue = odometervalue
	od.DateAdded = dateadded
	od.AddedBy = addedby
	return &od
}

// Prepare ...
func (o *Odometers) Prepare() {
	if o.DateAdded.IsZero() {
		o.DateAdded = time.Now()
	}
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	o.Units = "Kilometers"
}

// Validate ...
func (o *Odometers) Validate() error {

	if o.VehicleID == 0 {
		return errors.New("Odometers Required")
	}
	return nil
}

// Save ...
func (o *Odometers) Save(db *gorm.DB) (*Odometers, error) {
	if err := db.Debug().Model(o).Create(&o).Error; err != nil {
		return o, err
	}
	return o, nil
}

// GetOdometerByVehicleID ...
func (o *Odometers) GetOdometerByVehicleID(db *gorm.DB, vid uint64) (*[]Odometers, error) {
	var err error
	odometer := []Odometers{}
	if err = db.Debug().Model(o).Where("vehicle_id = ?", vid).Order("id DESC").Find(&odometer).Error; err != nil {
		return &odometer, err
	}

	return &odometer, nil
}
