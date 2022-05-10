package models

import "github.com/jinzhu/gorm"

// DashboardStats ...
type DashboardStats struct {
	MeterCount   int    `json:"meter_count"`
	GatewayCount int    `json:"gateway_count"`
	UserCount    int    `json:"user_count"`
	CompanyID    uint32 `json:"company_id"`
}

// GetUserDashboardStats ...
func (d *DashboardStats) GetUserDashboardStats(db *gorm.DB) *DashboardStats {
	// db.Debug().Model(&Meter{}).Where("id = ?", pid).Take(&p)
	db.Table("meters").Where("company_id = ?", d.CompanyID).Count(&d.MeterCount)
	db.Table("gateways").Where("company_id = ?", d.CompanyID).Count(&d.GatewayCount)
	db.Table("users").Where("company_id = ?", d.CompanyID).Count(&d.UserCount)
	return d
}
