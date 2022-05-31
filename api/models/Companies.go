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
	Status          uint8     `json:"status"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
	AddedBy         uint32    `json:"added_by"`
}

type CompanyShortDetails struct {
	ID          uint32 `json:"id"`
	CompanyName string `json:"company_name"`
}

// Prepare ...
func (p *Companies) Prepare() {
	p.ID = 0
	p.Status = 1
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

// Get ...
func (c *Companies) Get(db *gorm.DB) (*Companies, error) {
	err := db.Debug().Model(c).Where("id = ?", c.ID).First(&c).Error
	return c, err
}

// SaveCompanyDetails ...
func (c *Companies) Create(db *gorm.DB) (*Companies, error) {
	err := db.Debug().Model(c).Create(&c).Error
	return c, err
}

// Count Meters ...
func (c *Companies) Count(db *gorm.DB, roleid uint32) int {
	var count int
	query := db.Debug().Model(c)
	if roleid > 1002 {
		query = query.Where("added_by = ?", c.AddedBy)
	}

	query.Count(&count)
	return count
}

// List ...
func (c *Companies) List(db *gorm.DB, roleid uint32, offset, limit int) (*[]Companies, error) {
	var (
		err   error
		res   = []Companies{}
		query = db.Debug().Model(c)
	)

	if roleid > 1002 {
		query = query.Where("added_by = ?", c.AddedBy)
	}

	err = query.Offset(offset).Limit(limit).Find(&res).Error
	return &res, err
}

// Update Companies ...
func (c *Companies) Update(db *gorm.DB) (*Companies, error) {
	err := db.Debug().Model(c).Where("id = ?", c.ID).Updates(
		map[string]interface{}{
			"company_name":     c.CompanyName,
			"company_alias":    c.CompanyAlias,
			"company_phone":    c.CompanyPhone,
			"company_email":    c.CompanyEmail,
			"company_location": c.CompanyLocation,
			"status":           c.Status,
		},
	).Error

	return c, err
}

// Delete ...
func (c *Companies) Delete(db *gorm.DB) error {
	if err := db.Debug().Model(c).Where("id = ?", c.ID).Delete(&c).Error; err != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return errors.New("Company not found")
		}
		return db.Error
	}

	return nil
}
