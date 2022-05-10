package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Expenses ...
type Expenses struct {
	ID            uint32    `gorm:"primary_key;auto_increment" json:"id"`
	VehicleID     uint32    `gorm:"not null" json:"vehicle_id"`
	ServiceTypeID uint32    `gorm:"not null" json:"service_type_id"`
	TotalPrice    float64   `gorm:"not null" json:"total_price"`
	Description   string    `json:"description,omitempty"`
	ExpenseDate   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"expense_date"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	AddedBy       uint32    `gorm:"not null" json:"added_by"`
}

// InitializeExpenses ...
func InitializeExpenses(vehicleid, servicetypeid uint32, totalprice float64, description string, expensedate time.Time, addedby uint32) *Expenses {
	var ex Expenses
	ex.VehicleID = vehicleid
	ex.ServiceTypeID = servicetypeid
	ex.TotalPrice = totalprice
	ex.Description = description
	ex.ExpenseDate = expensedate
	ex.AddedBy = addedby
	return &ex
}

// Prepare ...
func (e *Expenses) Prepare() {
	if e.ExpenseDate.IsZero() {
		e.ExpenseDate = time.Now()
	}
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
}

// Validate ...
func (e *Expenses) Validate() error {

	if e.VehicleID == 0 {
		return errors.New("Vehicle Required")
	}
	return nil
}

// Save ...
func (e *Expenses) Save(db *gorm.DB) (*Expenses, error) {
	if err := db.Debug().Model(e).Create(&e).Error; err != nil {
		return e, err
	}
	return e, nil
}
