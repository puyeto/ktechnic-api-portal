package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Buildings struct {
	ID               uint32 `gorm:"primary_key;auto_increment" json:"id"`
	CompanyID        uint32 `gorm:"not null;" json:"company_id"`
	BuildingName     string `json:"building_name"`
	BuildingLocation string `json:"building_location"`
	Status           uint8  `json:"status"`
	AddedBy          uint32 `json:"added_by"`
}

// Validate ...
func (b *Buildings) Validate() error {
	if b.BuildingName == "" {
		return errors.New("Building Name is Required")
	}
	return nil
}

func (b *Buildings) Prepare() {
	b.Status = 1
}

// Get ...
func (b *Buildings) Get(db *gorm.DB) (*Buildings, error) {
	err := db.Debug().Model(b).Where("id = ?", b.ID).Take(&b).Error
	return b, err
}

// Count Meters ...
func (b *Buildings) Count(db *gorm.DB, roleid uint32) int {
	var count int
	query := db.Debug().Model(b)
	if roleid == 1002 {
		query = query.Where("company_id = ?", b.CompanyID)
	} else if roleid > 1002 {
		query = query.Where("added_by = ?", b.AddedBy)
	}

	query.Count(&count)
	return count
}

// List ...
func (b *Buildings) List(db *gorm.DB, roleid uint32, offset, limit int) (*[]Buildings, error) {
	var err error
	res := []Buildings{}

	query := db.Debug().Model(b)
	if roleid == 1002 {
		query = query.Where("company_id = ?", b.CompanyID)
	} else if roleid > 1002 {
		query = query.Where("added_by = ?", b.AddedBy)
	}

	err = query.Offset(offset).Limit(limit).Find(&res).Error
	return &res, err
}

// Create Building ...
func (b *Buildings) Create(db *gorm.DB) (*Buildings, error) {
	err := db.Debug().Model(b).Create(&b).Error
	return b, err
}

// Update Building ...
func (b *Buildings) Update(db *gorm.DB) (*Buildings, error) {
	err := db.Debug().Model(b).Where("id = ?", b.ID).Updates(
		map[string]interface{}{
			"building_name":     b.BuildingName,
			"building_location": b.BuildingLocation,
			"status":            b.Status,
		},
	).Error

	return b, err
}

// Delete ...
func (b *Buildings) Delete(db *gorm.DB) error {
	if err := db.Debug().Model(b).Where("id = ?", b.ID).Delete(&b).Error; err != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return errors.New("Building not found")
		}
		return db.Error
	}

	return nil
}
