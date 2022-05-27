package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Gateway ...
type Gateway struct {
	ID                 int                 `gorm:"primary_key;auto_increment" json:"id" db:"id"`
	SectionID          int                 `gorm:"null" json:"section_id"`
	CompanyID          uint32              `gorm:"not null;" json:"company_id"`
	GatewayName        string              `gorm:"not null;" json:"gateway_name" db:"gateway_name"`
	Status             int8                `gorm:"not null;" json:"status" db:"status"`
	GatewayDescription string              `gorm:"null;" json:"gateway_description"`
	CreatedAt          time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	AddedBy            uint32              `gorm:"-" json:"added_by"`
	Company            CompanyShortDetails `gorm:"-" json:"company_details"`
}

// Prepare ...
func (g *Gateway) Prepare() {
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	g.Status = 1
}

// Validate ...
func (g *Gateway) Validate() error {

	if g.GatewayName == "" {
		return errors.New("Gateway Name is Required")
	}

	return nil
}

// SaveGateway ...
func (g *Gateway) SaveGateway(db *gorm.DB) (*Gateway, error) {
	if err := db.Debug().Model(&Gateway{}).Create(&g).Error; err != nil {
		return &Gateway{}, err
	}

	return g, nil
}

// Count Meters ...
func (g *Gateway) CountGateways(db *gorm.DB, roleid uint32) int {
	var count int
	query := db.Debug().Model(g)
	if roleid == 1002 {
		query = query.Where("company_id = ?", g.CompanyID)
	} else if roleid > 1002 {
		query = query.Where("added_by = ?", g.AddedBy)
	}

	query.Count(&count)
	return count
}

// ListAllGateways ...
func (g *Gateway) ListAllGateways(db *gorm.DB, roleid, offset, limit uint32) (*[]Gateway, error) {
	var err error
	gateways := []Gateway{}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return &gateways, err
	}

	query := tx.Debug().Model(g)
	if roleid == 1002 {
		query = query.Where("company_id = ?", g.CompanyID)
	} else if roleid > 1002 {
		query = query.Where("added_by = ?", g.AddedBy)
	}

	err = query.Offset(offset).Limit(limit).Find(&gateways).Error

	if err != nil {
		return &gateways, err
	}

	if len(gateways) > 0 {
		for i := range gateways {
			tx.Debug().Model(&Companies{}).Where("id = ?", gateways[i].CompanyID).Take(&gateways[i].Company)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return &gateways, err
	}

	return &gateways, nil
}

// FindGatewayByID ...
func (g *Gateway) FindGatewayByID(db *gorm.DB, pid uint64) (*Gateway, error) {
	var err error
	err = db.Debug().Model(&Gateway{}).Where("id = ?", pid).Take(&g).Error
	if err != nil {
		return g, err
	}
	if g.ID != 0 {
		db.Debug().Table("companies").Model(&CompanyShortDetails{}).Where("id = ?", g.CompanyID).Take(&g.Company)
	}
	return g, nil
}

// UpdateAGateway ...
func (g *Gateway) UpdateAGateway(db *gorm.DB) (*Gateway, error) {

	var err error
	db.Debug().Model(&Gateway{}).Where("id = ?", g.ID).Take(&Gateway{}).UpdateColumns(
		map[string]interface{}{
			"gateway_name":        g.GatewayName,
			"section_id":          g.SectionID,
			"gateway_description": g.GatewayDescription,
			"updated_at":          g.UpdatedAt,
		},
	)
	err = db.Debug().Model(&Gateway{}).Where("id = ?", g.ID).Take(&g).Error
	if err != nil {
		return &Gateway{}, err
	}
	if g.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", g.CompanyID).Take(&g.Company).Error
		if err != nil {
			return &Gateway{}, err
		}
	}
	return g, nil
}

// DeleteAGateway ...
func (g *Gateway) DeleteAGateway(db *gorm.DB, vid uint32) error {
	if err := db.Debug().Model(&Gateway{}).Where("id = ?", vid).Delete(&Gateway{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return errors.New("Gateway not found")
		}
		return db.Error
	}

	return nil
}
