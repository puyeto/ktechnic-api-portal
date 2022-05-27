package models

// DashboardStats ...
type UserStats struct {
	MeterCount   int `json:"meter_count"`
	GatewayCount int `json:"gateway_count"`
	UserCount    int `json:"user_count"`
}

// // GetUserDashboardStats ...
// func (d *UserStats) GetUserDashboardStats(db *gorm.DB) *UserStats {
// 	meter := models.Meter{}

// 	// db.Debug().Model(&Meter{}).Where("id = ?", pid).Take(&p)
// 	db.Table("meters").Where("company_id = ?", d.CompanyID).Count(&d.MeterCount)
// 	db.Table("gateways").Where("company_id = ?", d.CompanyID).Count(&d.GatewayCount)
// 	db.Table("users").Where("company_id = ?", d.CompanyID).Count(&d.UserCount)
// 	return d
// }
