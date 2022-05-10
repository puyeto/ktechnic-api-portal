package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Companies ...
type Companies struct {
	ID              uint32    `json:"id"`
	CompanyName     string    `json:"company_name"`
	CompanyAlias    string    `json:"company_alias"`
	CompanyPhone    string    `json:"company_phone"`
	CompanyEmail    string    `json:"company_email"`
	CompanyLocation string    `json:"company_location"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
	AddedBy         uint32    `json:"added_by"`
}

// Prepare ...
func (p *Companies) Prepare() {
	p.ID = 0
	p.CompanyName = html.EscapeString(strings.ToUpper(p.CompanyName))
	p.CompanyAlias = html.EscapeString(strings.ToUpper(p.CompanyAlias))
	p.CompanyLocation = html.EscapeString(strings.TrimSpace(p.CompanyLocation))
	p.CompanyEmail = html.EscapeString(strings.ToLower(p.CompanyEmail))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate ...
func (p *Companies) Validate() error {

	if p.CompanyName == "" {
		return errors.New("Required Company Name")
	}
	if p.CompanyAlias == "" {
		return errors.New("Required Company Alias Name")
	}
	if p.CompanyPhone == "" {
		return errors.New("Required Company Phone No")
	}
	if p.CompanyEmail == "" {
		return errors.New("Required Company Email")
	}
	if p.CompanyLocation == "" {
		return errors.New("Required Company Location")
	}
	return nil
}

// SaveCompanyDetails ...
func (p *Companies) SaveCompanyDetails(db *gorm.DB) (*Companies, error) {
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return &Companies{}, err
	}

	err := tx.Debug().Model(&Companies{}).Create(&p).Error
	if err != nil {
		return &Companies{}, err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return &Companies{}, err
	}

	return p, nil
}

// List ...
func (p *Companies) List(db *gorm.DB) (*[]Companies, error) {
	var err error
	con := []Companies{}
	err = db.Debug().Model(&Companies{}).Limit(100).Find(&con).Error
	if err != nil {
		return &[]Companies{}, err
	}

	return &con, nil
}
