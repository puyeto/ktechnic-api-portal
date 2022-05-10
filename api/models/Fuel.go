package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Fuel ...
type Fuel struct {
	ID               uint32    `gorm:"primary_key;auto_increment" json:"id"`
	VehicleID        uint32    `gorm:"not null" json:"vehicle_id"`
	DriverID         uint32    `gorm:"not null" json:"driver_id"`
	OdometerID       uint32    `gorm:"not null" json:"odometer_id"`
	OdometerValue    float64   `json:"odometer_value,omitempty"`
	ExpenseID        uint32    `gorm:"not null" json:"expense_id"`
	Litres           float64   `gorm:"not null" json:"litres"`
	PricePerLitre    float64   `gorm:"not null" json:"price_per_litre"`
	TotalPrice       float64   `gorm:"not null" json:"total_price"`
	InvoiceReference string    `json:"invoice_reference"`
	VendorID         uint32    `gorm:"not null" json:"vendor_id"`
	Description      string    `json:"description,omitempty"`
	DateFueled       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date_fueled"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	AddedBy          uint32    `gorm:"not null" json:"added_by"`
}

// Prepare ...
func (f *Fuel) Prepare() {
	if f.DateFueled.IsZero() {
		f.DateFueled = time.Now()
	}
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	f.TotalPrice = f.Litres * f.PricePerLitre
}

// Validate ...
func (f *Fuel) Validate() error {

	if f.VehicleID == 0 {
		return errors.New("Vehicle Required")
	}
	if f.DriverID == 0 {
		return errors.New("Driver Required")
	}
	if f.Litres == 0 {
		return errors.New("Litres is Required")
	}
	if f.PricePerLitre == 0 {
		return errors.New("Price Per Litre is Required")
	}
	return nil
}

// GetRefuelsByVehicleID ...
func (f *Fuel) GetRefuelsByVehicleID(db *gorm.DB, vid uint64) (*[]Fuel, error) {
	var err error
	fuel := []Fuel{}
	if err = db.Debug().Model(f).Where("vehicle_id = ?", vid).Order("id DESC").Find(&fuel).Error; err != nil {
		return &fuel, err
	}

	return &fuel, nil
}

// Save Refuels ...
func (f *Fuel) Save(db *gorm.DB) (*Fuel, error) {
	if err := db.Debug().Model(f).Create(&f).Error; err != nil {
		return f, err
	}
	return f, nil
}
