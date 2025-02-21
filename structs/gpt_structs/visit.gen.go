package gpt_structs

import (
	"time"
)

const TableNameVisit = "visit"

// Visit mapped from table <visit>
type Visit struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID     *int64    `gorm:"column:customer_id" json:"customer_id"`
	UserIDExecutor *int32    `gorm:"column:user_id_executor" json:"user_id_executor"`
	UserIDHolder   *int32    `gorm:"column:user_id_holder" json:"user_id_holder"`
	CheckinAt      time.Time `gorm:"column:checkin_at" json:"checkin_at"`
	CheckoutAt     time.Time `gorm:"column:checkout_at" json:"checkout_at"`
	VisitStatus    *string   `gorm:"column:visit_status" json:"visit_status"`
	IsFakeGps      *int16    `gorm:"column:is_fake_gps" json:"is_fake_gps"`
	ProvinceID     *int16    `gorm:"column:province_id" json:"province_id"`
	RegencyID      *int16    `gorm:"column:regency_id" json:"regency_id"`
	DistrictID     *int16    `gorm:"column:district_id" json:"district_id"`
	SubDistrictID  *int32    `gorm:"column:sub_district_id" json:"sub_district_id"`
	SrID           *int16    `gorm:"column:sr_id" json:"sr_id"`
	RayonID        *int16    `gorm:"column:rayon_id" json:"rayon_id"`
	BranchID       *int32    `gorm:"column:branch_id" json:"branch_id"`
	Latitude       *string   `gorm:"column:latitude" json:"latitude"`
	Longitude      *string   `gorm:"column:longitude" json:"longitude"`
}

// TableName Visit's table name
func (*Visit) TableName() string {
	return TableNameVisit
}
