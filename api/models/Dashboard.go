package models

// DashboardStats ...
type UserStats struct {
	TotalConsumption float32 `json:"total_consumption"`
	TotalSpent       float32 `json:"total_spent"`
	MeterCount       int     `json:"meter_count"`
	OpenValve        int     `json:"open_valve"`
	ClosedValve      int     `json:"closed_valve"`
	GatewayCount     int     `json:"gateway_count"`
	UserCount        int     `json:"user_count"`
}

type MaterStats struct {
	MeterNumber               int64   `json:"meter_number"`
	MeterUnitsBalance         float32 `json:"meter_units_balance"`
	MeterWalletBalance        float32 `json:"meter_wallet_balance"`
	MeterConsumptionThisWeek  float32 `json:"meter_consumption_this_week"`
	MeterConsumptionThisMonth float32 `json:"meter_consumption_this_month"`
	MeterLastSeen             string  `json:"meter_last_seen"`
	MeterValveStatus          bool    `json:"meter_valve_status"`
	MeterDetails              *Meter  `json:"meter_details"`
}
