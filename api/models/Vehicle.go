package models

import (
	"errors"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Vehicle ...
type Vehicle struct {
	ID              uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Driver          User      `json:"driver"`
	DriverID        uint32    `gorm:"not null" json:"driver_id"`
	OwnerID         uint32    `gorm:"not null" json:"owner_id"`
	VehicleOwner    User      `json:"vehicle_owner"`
	CompanyID       uint32    `gorm:"not null;" json:"company_id"`
	RegistrationNo  string    `gorm:"size:55;not null;" json:"registration_no"`
	VehicleStringID string    `gorm:"size:55;not null;unique" json:"vehicle_string_id"`
	Make            string    `gorm:"size:255;not null;" json:"make"`
	Model           string    `gorm:"size:255;not null;" json:"model"`
	ChassisNo       string    `gorm:"size:255;not null" json:"chassis_no"`
	NoOfPassengers  int       `gorm:"not null" json:"no_of_passengers"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	AddedBy         uint32    `gorm:"not null" json:"added_by"`
}

// Prepare ...
func (p *Vehicle) Prepare() {
	p.Driver = User{}
	p.RegistrationNo = html.EscapeString(strings.TrimSpace(strings.ToUpper(p.RegistrationNo)))
	p.Make = html.EscapeString(strings.TrimSpace(strings.ToUpper(p.Make)))
	p.Model = html.EscapeString(strings.TrimSpace(strings.ToUpper(p.Model)))
	p.ChassisNo = html.EscapeString(strings.TrimSpace(strings.ToUpper(p.ChassisNo)))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	p.VehicleStringID = strings.ToLower(reg.ReplaceAllString(p.RegistrationNo, "_"))
}

// Validate ...
func (p *Vehicle) Validate() error {

	if p.RegistrationNo == "" {
		return errors.New("Vehicle Registration No. is Required")
	}
	if p.ChassisNo == "" {
		return errors.New("Vehicle Chassis No. is Required")
	}
	if p.Make == "" {
		return errors.New("Vehicle Make is Required")
	}
	if p.Model == "" {
		return errors.New("Vehicle Model is Required")
	}
	if p.DriverID < 1 {
		return errors.New("Vehicle Driver is Required")
	}
	return nil
}

// SaveVehicle ...
func (p *Vehicle) SaveVehicle(db *gorm.DB) (*Vehicle, error) {
	var err error
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return &Vehicle{}, err
	}

	if err = tx.Debug().Model(&Vehicle{}).Create(&p).Error; err != nil {
		tx.Rollback()
		return &Vehicle{}, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &Vehicle{}, err
	}

	return p, nil
}

// ListAllVehicles ...
func (p *Vehicle) ListAllVehicles(db *gorm.DB) (*[]Vehicle, error) {
	var err error
	vehicles := []Vehicle{}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return &[]Vehicle{}, err
	}

	if p.CompanyID > 0 {
		err = tx.Debug().Where("company_id = ?", p.CompanyID).Model(&Vehicle{}).Limit(100).Find(&vehicles).Error
	} else {
		err = tx.Debug().Model(&Vehicle{}).Limit(100).Find(&vehicles).Error
	}

	if err != nil {
		tx.Rollback()
		return &[]Vehicle{}, err
	}

	if len(vehicles) > 0 {
		for i := range vehicles {
			if err := tx.Debug().Model(&User{}).Where("id = ?", vehicles[i].DriverID).Take(&vehicles[i].Driver).Error; err != nil {
				return &[]Vehicle{}, err
			}
			if err := tx.Debug().Model(&User{}).Where("id = ?", vehicles[i].OwnerID).Take(&vehicles[i].VehicleOwner).Error; err != nil {
				return &[]Vehicle{}, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &[]Vehicle{}, err
	}

	return &vehicles, nil
}

// FindVehicleByID ...
func (p *Vehicle) FindVehicleByID(db *gorm.DB, pid uint64) (*Vehicle, error) {
	var err error
	err = db.Debug().Model(&Vehicle{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Vehicle{}, err
	}
	if p.ID != 0 {
		if err = db.Debug().Model(&User{}).Where("id = ?", p.DriverID).Take(&p.Driver).Error; err != nil {
			return &Vehicle{}, err
		}

		if err = db.Debug().Model(&User{}).Where("id = ?", p.OwnerID).Take(&p.VehicleOwner).Error; err != nil {
			return &Vehicle{}, err
		}
	}
	return p, nil
}

// FindVehicleByQueryString ...
func (p *Vehicle) FindVehicleByQueryString(db *gorm.DB) (*Vehicle, error) {
	var err error
	err = db.Debug().Model(&Vehicle{}).Where("vehicle_string_id = ?", p.VehicleStringID).Take(&p).Error
	if err != nil {
		return &Vehicle{}, err
	}

	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.DriverID).Take(&p.Driver).Error
		if err != nil {
			return &Vehicle{}, err
		}

		if err := db.Debug().Model(&User{}).Where("id = ?", p.OwnerID).Take(&p.VehicleOwner).Error; err != nil {
			return &Vehicle{}, err
		}
	}
	return p, nil
}

// UpdateAVehicle ...
func (p *Vehicle) UpdateAVehicle(db *gorm.DB) (*Vehicle, error) {

	var err error
	db.Debug().Model(&Vehicle{}).Where("id = ?", p.ID).Take(&Vehicle{}).UpdateColumns(
		map[string]interface{}{
			"registration_no":   p.RegistrationNo,
			"make":              p.Make,
			"model":             p.Model,
			"chassis_no":        p.ChassisNo,
			"updated_at":        time.Now(),
			"no_of_passengers":  p.NoOfPassengers,
			"vehicle_string_id": p.VehicleStringID,
		},
	)
	// err = db.Debug().Model(&Vehicle{}).Where("id = ?", p.ID).Take(&p).Error
	if err != nil {
		return &Vehicle{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.DriverID).Take(&p.Driver).Error
		if err != nil {
			return &Vehicle{}, err
		}
	}
	return p, nil
}

// DeleteAVehicle ...
func (p *Vehicle) DeleteAVehicle(db *gorm.DB, vid uint32) (int64, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return 0, err
	}

	err := tx.Debug().Model(&Vehicle{}).Where("id = ?", vid).Take(&p).Error
	if err != nil {
		return 0, err
	}

	if err = tx.Debug().Model(&Vehicle{}).Where("id = ?", vid).Delete(&Vehicle{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Vehicle not found")
		}
		return 0, db.Error
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return db.RowsAffected, nil
}
