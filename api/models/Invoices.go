package models

import (
	"errors"
	"fmt"
	"math"

	"github.com/jinzhu/gorm"
)

type CalculateTotalsDetails struct {
	MeterID int32   `json:"meter_id"`
	Amount  float32 `json:"amount_paid"`
}

func (c *CalculateTotalsDetails) Validate() error {
	if c.MeterID == 0 {
		return errors.New("Meter details is required")
	}

	if c.Amount == 0 {
		return errors.New("Amount is required")
	}
	return nil
}

type CalculateTotals struct {
	MeterID       int32   `json:"meter_id"`
	MeterName     string  `json:"meter_name"`
	MeterNumber   string  `json:"meter_number"`
	PricePlanID   int32   `json:"price_plan_id"`
	PricePlanName string  `json:"price_plan_name"`
	Amount        float32 `json:"total_amount_paid"`
	Deductions    float32 `json:"total_deductions"`
	UnitsAmount   float32 `json:"total_units_amount"`
	AmountPerUnit float32 `json:"amount_per_unit"`
	Units         float64 `json:"total_units"`
}

func (c *CalculateTotalsDetails) CalculateTotal(db *gorm.DB) (*CalculateTotals, error) {
	var calc CalculateTotals
	calc.MeterID = c.MeterID
	calc.Amount = c.Amount

	tx := db.Begin().Debug()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	// get meter Details
	if err := tx.Table("meters").Select("meter_name, meter_number, price_plan_id").Where("id = ?", c.MeterID).Take(&calc).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// get meter Price Plan
	fmt.Println(calc.PricePlanID)
	if err := tx.Table("price_plans").Select("price_plan_name, amount_per_unit").Where("id = ?", calc.PricePlanID).Take(&calc).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// get deductions
	var price []PricePlanDetails
	db.Debug().Model(&PricePlanDetails{}).Where("price_plan_id = ?", calc.PricePlanID).Order("price_detail_type ASC").Find(&price)
	for i := range price {
		if price[i].PriceDetailType == "fixed" {
			calc.Deductions = calc.Deductions + float32(price[i].PriceDetailValue)
		}

		if price[i].PriceDetailType == "percent" {
			calc.Deductions = calc.Deductions + (float32(price[i].PriceDetailValue) / 100 * calc.Amount)
		}
	}
	calc.UnitsAmount = calc.Amount - calc.Deductions

	// calculate Units
	calc.Units = math.Round(float64(calc.UnitsAmount/calc.AmountPerUnit)*100) / 100

	return &calc, tx.Commit().Error
}

func (c *CalculateTotalsDetails) MakePayment(db *gorm.DB) (*CalculateTotals, error) {
	var calc CalculateTotals
	calc.MeterID = c.MeterID
	calc.Amount = c.Amount

	calctotals, err := c.CalculateTotal(db)
	if err != nil {
		return nil, err
	}

	// Create invoice
	if err := db.Debug().Table("invoices").Create(&calctotals).Error; err != nil {
		return &calc, err
	}

	// create Payment

	return &calc, nil
}
