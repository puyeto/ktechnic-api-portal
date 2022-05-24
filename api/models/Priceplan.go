package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type PricePlan struct {
	ID               uint32             `json:"id" gorm:"primary_key;auto_increment"`
	CompanyID        uint32             `json:"company_id" gorm:"not null"`
	PricePlanName    string             `json:"price_plan_name" gorm:"not null"`
	AmountPerUnit    float64            `json:"amount_per_unit" gorm:"not null"`
	AddedBy          uint32             `json:"added_by"`
	PricePlanDetails []PricePlanDetails `json:"price_plan_details"`
}

type PricePlanDetails struct {
	ID               uint32  `json:"id" gorm:"id"`
	PricePlanID      uint32  `json:"price_plan_id" gorm:"not null"`
	PriceDetailName  string  `json:"price_detail_name" gorm:"not null"`
	PriceDetailValue float64 `json:"price_detail_value" gorm:"default:0.00"`
	PriceDetailType  string  `json:"price_detail_type"`
}

// Validate ...
func (p *PricePlan) Validate() error {
	if p.PricePlanName == "" {
		return errors.New("Plan Name is Required")
	}
	if p.AmountPerUnit == 0 {
		return errors.New("Amount is Required")
	}
	return nil
}

func (p *PricePlan) ListPricePlan(db *gorm.DB, roleid uint32) (*[]PricePlan, error) {
	plans := []PricePlan{}

	query := db.Debug().Model(p)
	if roleid == 1002 {
		query = query.Where("company_id = ?", p.CompanyID)
	} else if roleid > 1002 {
		query = query.Where("added_by = ?", p.AddedBy)
	}

	if err := query.Limit(100).Find(&plans).Error; err != nil {
		return &plans, err
	}

	if len(plans) > 0 {
		for i := range plans {
			db.Debug().Model(&PricePlanDetails{}).Where("price_plan_id = ?", plans[i].ID).Find(&plans[i].PricePlanDetails)
		}
	}

	return &plans, nil
}

func (p *PricePlan) SavePricePlan(db *gorm.DB) (*PricePlan, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return p, err
	}

	if err := tx.Debug().Model(p).Create(&p).Error; err != nil {
		return p, err
	}

	if err := tx.Commit().Error; err != nil {
		return p, err
	}

	return p, nil
}

func (p *PricePlan) DeletePricePlan(db *gorm.DB, pid uint32) error {
	if err := db.Debug().Model(p).Where("id = ?", pid).Delete(p).Error; err != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return errors.New("Plan not found")
		}
		return db.Error
	}

	db.Debug().Model(&PricePlanDetails{}).Where("price_plan_id = ?", pid).Delete(&PricePlanDetails{})

	return nil
}
