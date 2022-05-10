package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Gateway ...
type Gateway struct {
	ID                 int       `gorm:"primary_key;auto_increment" json:"id" db:"id"`
	SectionID          int       `gorm:"null" json:"section_id"`
	CompanyID          uint32    `gorm:"not null;" json:"company_id"`
	Company            Companies `json:"company_details"`
	GatewayName        string    `gorm:"not null;" json:"gateway_name" db:"gateway_name"`
	Status             int8      `gorm:"not null;" json:"status" db:"status"`
	GatewayDescription string    `gorm:"null;" json:"gateway_description"`
	CreatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	AddedBy            uint32    `gorm:"not null" json:"added_by"`
}

// Prepare ...
func (p *Gateway) Prepare() {
	p.Company = Companies{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.Status = 1
}

// Validate ...
func (p *Gateway) Validate() error {

	if p.GatewayName == "" {
		return errors.New("Gateway Name is Required")
	}

	return nil
}

// SaveGateway ...
func (p *Gateway) SaveGateway(db *gorm.DB) (*Gateway, error) {
	var err error
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return &Gateway{}, err
	}

	if err = tx.Debug().Model(&Gateway{}).Create(&p).Error; err != nil {
		tx.Rollback()
		return &Gateway{}, err
	}
	if p.ID > 0 {
		if err = tx.Debug().Model(&User{}).Where("id = ?", p.AddedBy).UpdateColumn("gateway_count", gorm.Expr("gateway_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			return &Gateway{}, err
		}

		if err = tx.Debug().Model(&Companies{}).Where("id = ?", p.CompanyID).UpdateColumn("gateway_count", gorm.Expr("gateway_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			return &Gateway{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &Gateway{}, err
	}

	return p, nil
}

// ListAllGateways ...
func (p *Gateway) ListAllGateways(db *gorm.DB) (*[]Gateway, error) {
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

	if p.CompanyID > 0 {
		err = tx.Debug().Where("company_id = ?", p.CompanyID).Model(&Gateway{}).Limit(100).Find(&gateways).Error
	} else {
		err = tx.Debug().Model(&Gateway{}).Limit(100).Find(&gateways).Error
	}

	if err != nil {
		tx.Rollback()
		return &gateways, err
	}

	if len(gateways) > 0 {
		for i := range gateways {
			tx.Debug().Model(&Companies{}).Where("id = ?", gateways[i].CompanyID).Take(&gateways[i].Company)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return &gateways, err
	}

	return &gateways, nil
}

// FindGatewayByID ...
func (p *Gateway) FindGatewayByID(db *gorm.DB, pid uint64) (*Gateway, error) {
	var err error
	err = db.Debug().Model(&Gateway{}).Where("id = ?", pid).Take(&p).Error
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

// UpdateAGateway ...
func (p *Gateway) UpdateAGateway(db *gorm.DB) (*Gateway, error) {

	var err error
	db.Debug().Model(&Gateway{}).Where("id = ?", p.ID).Take(&Gateway{}).UpdateColumns(
		map[string]interface{}{
			"gateway_name":        p.GatewayName,
			"section_id":          p.SectionID,
			"gateway_description": p.GatewayDescription,
			"updated_at":          p.UpdatedAt,
		},
	)
	err = db.Debug().Model(&Gateway{}).Where("id = ?", p.ID).Take(&p).Error
	if err != nil {
		return &Gateway{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.CompanyID).Take(&p.Company).Error
		if err != nil {
			return &Gateway{}, err
		}
	}
	return p, nil
}

// DeleteAGateway ...
func (p *Gateway) DeleteAGateway(db *gorm.DB, vid uint32) (int64, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return 0, err
	}

	err := tx.Debug().Model(&Gateway{}).Where("id = ?", vid).Take(&p).Error
	if err != nil {
		return 0, err
	}

	if err = tx.Debug().Model(&User{}).Where("id = ?", p.AddedBy).UpdateColumn("gateway_count", gorm.Expr("gateway_count - ?", 1)).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = tx.Debug().Model(&Companies{}).Where("id = ?", p.CompanyID).UpdateColumn("gateway_count", gorm.Expr("gateway_count - ?", 1)).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = tx.Debug().Model(&Gateway{}).Where("id = ?", vid).Delete(&Gateway{}).Error; err != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Gateway not found")
		}
		return 0, db.Error
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return db.RowsAffected, nil
}
