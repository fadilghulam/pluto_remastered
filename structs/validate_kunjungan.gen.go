// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package structs

import (
	"time"
)

const TableNameValidateKunjungan = "validate_kunjungan"

// ValidateKunjungan mapped from table <validate_kunjungan>
type ValidateKunjungan struct {
	ID                int64     `gorm:"column:id;primaryKey;default:null" json:"id"`
	CustomerID        int64     `gorm:"column:customer_id;not null" json:"customer_id"`
	KunjunganID       int64     `gorm:"column:kunjungan_id;not null" json:"kunjungan_id"`
	UserID            int32     `gorm:"column:user_id;not null" json:"user_id"`
	Datetime          time.Time `gorm:"column:datetime;not null;default:now()" json:"datetime"`
	IsExist           int16     `gorm:"column:is_exist;not null" json:"is_exist"`
	ExistInfo         string    `gorm:"column:exist_info;default:null" json:"exist_info"`
	IsValid           int16     `gorm:"column:is_valid;default:null" json:"is_valid"`
	ValidateInfo      string    `gorm:"column:validate_info;default:null" json:"validate_info"`
	CreatedAt         time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	LatitudeLongitude string    `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
}

// TableName ValidateKunjungan's table name
func (*ValidateKunjungan) TableName() string {
	return TableNameValidateKunjungan
}
