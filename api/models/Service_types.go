package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

// ServiceTypes ...
type ServiceTypes struct {
	ID          uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name        string `gorm:"size:55;not null" json:"name"`
	Description string `gorm:"size:55;not null" json:"description"`
	AddedBy     uint32 `gorm:"not null" json:"added_by"`
}

// Prepare ...
func (s *ServiceTypes) Prepare() {
	s.Name = html.EscapeString(strings.TrimSpace(s.Name))
	if s.Description == "" {
		s.Description = s.Name
	}
}

// Validate ...
func (s *ServiceTypes) Validate() error {

	if s.Name == "" {
		return errors.New("Required Name")
	}
	return nil
}

// List List Service Types ...
func (s *ServiceTypes) List(db *gorm.DB) (*[]ServiceTypes, error) {
	var err error
	con := []ServiceTypes{}
	err = db.Debug().Model(&ServiceTypes{}).Order("name ASC").Limit(100).Find(&con).Error
	if err != nil {
		return &con, err
	}

	return &con, nil
}

// Save Service Types ...
func (s *ServiceTypes) Save(db *gorm.DB) (*ServiceTypes, error) {
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return &ServiceTypes{}, err
	}

	err := tx.Debug().Model(&ServiceTypes{}).Create(&s).Error
	if err != nil {
		return &ServiceTypes{}, err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return &ServiceTypes{}, err
	}

	return s, nil
}
