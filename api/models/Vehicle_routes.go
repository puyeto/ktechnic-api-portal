package models

import (
	"errors"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// VehicleRoutes routes structure
type VehicleRoutes struct {
	ID           uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Departure    string    `gorm:"size:100;not null;" json:"departure"`
	Destination  string    `gorm:"size:100;not null;" json:"destination"`
	DepDes       string    `gorm:"size:100;not null;" json:"depdes"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdatedBy    uint32    `json:"updated_by"`
	CompanyID    uint32    `json:"company_id"`
	DefaultPrice float64   `json:"default_price"`
}

// Prepare ...
func (u *VehicleRoutes) Prepare() {
	u.ID = 0
	u.Departure = html.EscapeString(strings.TrimSpace(u.Departure))
	u.Destination = html.EscapeString(strings.TrimSpace(u.Destination))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// Validate ...
func (u *VehicleRoutes) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Departure == "" {
			return errors.New("Required Departure")
		}
		if u.Destination == "" {
			return errors.New("Required Destination")
		}
		return nil
	default:
		if u.Departure == "" {
			return errors.New("Required Departure")
		}
		if u.Destination == "" {
			return errors.New("Required Destination")
		}
		return nil
	}
}

// SaveVehicleRoute ...
func (u *VehicleRoutes) SaveVehicleRoute(db *gorm.DB) (*VehicleRoutes, error) {
	depdesfieldname := u.Departure + " " + u.Destination
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return &VehicleRoutes{}, err
	}
	u.DepDes = strings.ToLower(reg.ReplaceAllString(depdesfieldname, "_"))

	if err = db.Debug().Create(&u).Error; err != nil {
		return &VehicleRoutes{}, err
	}
	return u, nil
}

// ListVehicleRoutes ...
func (u *VehicleRoutes) ListVehicleRoutes(db *gorm.DB) (*[]VehicleRoutes, error) {
	var err error
	routes := []VehicleRoutes{}

	if u.CompanyID > 0 {
		err = db.Debug().Where("company_id = ?", u.CompanyID).Model(&VehicleRoutes{}).Limit(100).Find(&routes).Error
	} else {
		err = db.Debug().Model(&VehicleRoutes{}).Limit(100).Find(&routes).Error
	}

	if err != nil {
		return &[]VehicleRoutes{}, err
	}
	return &routes, err
}
