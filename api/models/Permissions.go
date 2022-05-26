package models

import (
	"errors"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Permissions ...
type Permissions struct {
	ID           uint32    `gorm:"primary_key;auto_increment" json:"id"`
	FieldName    string    `gorm:"size:55;not null" json:"field_name"`
	ColumnName   string    `gorm:"size:55;not null" json:"column_name"`
	DefaultValue bool      `gorm:"-" json:"default_value"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	AddedBy      uint32    `gorm:"not null" json:"-"`
}

// Prepare ...
func (p *Permissions) Prepare() {
	p.ID = 0
	p.FieldName = html.EscapeString(strings.TrimSpace(p.FieldName))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate ...
func (p *Permissions) Validate() error {

	if p.FieldName == "" {
		return errors.New("Required Field Name")
	}
	return nil
}

// List List permissions ...
func (p *Permissions) List(db *gorm.DB) (*[]Permissions, error) {
	var err error
	con := []Permissions{}
	err = db.Debug().Model(&Permissions{}).Limit(100).Find(&con).Error
	if err != nil {
		return &con, err
	}

	return &con, nil
}

// Save Permissions ...
func (p *Permissions) Save(db *gorm.DB) (*Permissions, error) {
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return &Permissions{}, err
	}

	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return &Permissions{}, err
	}
	p.ColumnName = strings.ToLower(reg.ReplaceAllString(p.FieldName, "_"))
	// defaultvalue := fmt.Sprintf("%.2f", p.DefaultValue)

	query := "ALTER TABLE `user_permissions` ADD `" + p.ColumnName + "` TINYINT(1) NOT NULL DEFAULT '0';"
	tx.Debug().Exec(query)

	err = tx.Debug().Model(&Permissions{}).Create(&p).Error
	if err != nil {
		return &Permissions{}, err
	}

	if err = tx.Commit().Error; err != nil {
		return &Permissions{}, err
	}

	return p, nil
}

// Roles ...
type Roles struct {
	ID              uint32 `gorm:"primary_key;auto_increment" json:"id"`
	RoleName        string `gorm:"size:55;not null" json:"role_name"`
	RoleDescription string `gorm:"size:55;not null" json:"role_description"`
}

// List Roles ...
func (p *Roles) List(db *gorm.DB, cid, rid uint32) (*[]Roles, error) {
	var err error
	con := []Roles{}

	whre := "id >= ?"
	if rid == 1001 {
		whre = "id = ?"
	}

	err = db.Debug().Where(whre, rid).Model(&Roles{}).Limit(10).Find(&con).Error
	if err != nil {
		return &con, err
	}

	return &con, nil
}
