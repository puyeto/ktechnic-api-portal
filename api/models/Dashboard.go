package models

// DashboardStats ...
type UserStats struct {
	MeterCount   int `json:"meter_count"`
	GatewayCount int `json:"gateway_count"`
	UserCount    int `json:"user_count"`
}

type MaterStats struct {
	MeterNumber               string  `json:"meter_number"`
	MeterUnitsBalance         int     `json:"meter_units_balance"`
	MeterWalletBalance        float32 `json:"meter_wallet_balance"`
	MeterConsumptionThisWeek  float32 `json:"meter_consumption_this_week"`
	MeterConsumptionThisMonth float32 `json:"meter_consumption_this_month"`
	MeterLastSeen             string  `json:"meter_last_seen"`
	MeterDetails              *Meter  `json:"meter_details"`
}
