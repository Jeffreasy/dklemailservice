package models

// Distance (voorheen route_funds) definieert de beschikbare routes.
type Distance struct {
	Route           string `json:"route" gorm:"primaryKey"`
	FundAmount      int    `json:"fund_amount" gorm:"column:fund_amount"` // V26_08: Renamed from 'amount'
	RegistrationFee int    `json:"registration_fee"`
}

// TableName specificeert de tabelnaam voor GORM
func (Distance) TableName() string {
	return "distances"
}
