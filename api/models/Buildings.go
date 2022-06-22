package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Buildings struct {
	ID                   uint32                 `gorm:"primary_key;auto_increment" json:"id"`
	CompanyID            uint32                 `gorm:"not null;" json:"company_id"`
	GatewayID            uint32                 `gorm:"not null;" json:"gateway_id"`
	BuildingName         string                 `json:"building_name"`
	BuildingLocation     string                 `json:"building_location"`
	Status               uint8                  `json:"status"`
	AddedBy              uint32                 `json:"added_by"`
	BuildingHouseNumbers []BuildingHouseNumbers `gorm:"-" json:"building_house_numbers"`
}

// Validate ...
func (b *Buildings) Validate() error {
	if b.BuildingName == "" {
		return errors.New("Building Name is Required")
	}
	if b.GatewayID == 0 {
		return errors.New("Building gateway is Required")
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

type BuildingWithHouseNo struct {
	Buildings
	BuildingHouseNumbers []BuildingHouseNumbers `gorm:"-" json:"building_house_numbers"`
}

// Get ...
func (b *BuildingWithHouseNo) GetWithHouseNumbers(db *gorm.DB) (*BuildingWithHouseNo, error) {
	if err := db.Debug().Table("buildings").Model(b).Where("id = ?", b.ID).Take(&b).Error; err != nil {
		return b, err
	}

	err := db.Debug().Model(b).Where("building_id = ?", b.ID).Limit(100).Find(&b.BuildingHouseNumbers).Error
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

type BuildingHouseNumber struct {
	BuildingID           uint32                 `json:"building_id"`
	BuildingHouseNumbers []BuildingHouseNumbers `json:"building_house_numbers"`
}

type BuildingHouseNumbers struct {
	ID          uint32 `gorm:"primary_key;auto_increment" json:"id"`
	BuildingID  uint32 `gorm:"->:false;<-:create" json:"building_id"`
	HouseDetail string `json:"house_detail"`
}

// Create Building ...
func (b *BuildingHouseNumber) Create(db *gorm.DB) error {
	var err error
	for _, val := range b.BuildingHouseNumbers {
		if val.HouseDetail != "" {
			val.BuildingID = b.BuildingID
			if err = db.Debug().Model(BuildingHouseNumbers{}).Create(&val).Error; err != nil {
				break
			}
		}
	}

	return err
}
