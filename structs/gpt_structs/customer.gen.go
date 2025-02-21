package gpt_structs

import (
	"time"
)

const TableNameCustomer = "customer"

// Customer mapped from table <customer>
type Customer struct {
	ID                int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	OwnerName         *string   `gorm:"column:owner_name" json:"owner_name"`
	OutletName        *string   `gorm:"column:outlet_name" json:"outlet_name"`
	CustomerTypeID    *int16    `gorm:"column:customer_type_id" json:"customer_type_id"`
	UserIDHolder      *int32    `gorm:"column:user_id_holder" json:"user_id_holder"`
	RegisteredAt      time.Time `gorm:"column:registered_at" json:"registered_at"`
	UserIDRegistrant  *int32    `gorm:"column:user_id_registrant" json:"user_id_registrant"`
	IsVerified        *int16    `gorm:"column:is_verified" json:"is_verified"`
	VerifiedAt        time.Time `gorm:"column:verified_at" json:"verified_at"`
	UserIDVerificator *int32    `gorm:"column:user_id_verificator" json:"user_id_verificator"`
	ProvinceID        *int16    `gorm:"column:province_id" json:"province_id"`
	RegencyID         *int16    `gorm:"column:regency_id" json:"regency_id"`
	DistrictID        *int16    `gorm:"column:district_id" json:"district_id"`
	SubDistrictID     *int32    `gorm:"column:sub_district_id" json:"sub_district_id"`
	SrID              *int16    `gorm:"column:sr_id" json:"sr_id"`
	RayonID           *int16    `gorm:"column:rayon_id" json:"rayon_id"`
	BranchID          *int32    `gorm:"column:branch_id" json:"branch_id"`
	Latitude          *string   `gorm:"column:latitude" json:"latitude"`
	Longitude         *string   `gorm:"column:longitude" json:"longitude"`
}

// TableName Customer's table name
func (*Customer) TableName() string {
	return TableNameCustomer
}
