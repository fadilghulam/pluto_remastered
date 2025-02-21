package gpt_structs

import (
	"time"
)

const TableNameTransaction = "transaction"

// Transaction mapped from table <transaction>
type Transaction struct {
	ID                  int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	TransactionAt       time.Time `gorm:"column:transaction_at" json:"transaction_at"`
	TransactionTypeID   *int16    `gorm:"column:transaction_type_id" json:"transaction_type_id"`
	UserIDExecutor      *int32    `gorm:"column:user_id_executor" json:"user_id_executor"`
	UserIDHolder        *int32    `gorm:"column:user_id_holder" json:"user_id_holder"`
	CustomerID          *int64    `gorm:"column:customer_id" json:"customer_id"`
	GrandTotal          *int64    `gorm:"column:grand_total" json:"grand_total"`
	RemainingReceivable *int64    `gorm:"column:remaining_receivable" json:"remaining_receivable"`
	PaidAt              time.Time `gorm:"column:paid_at" json:"paid_at"`
	ProvinceID          *int16    `gorm:"column:province_id" json:"province_id"`
	RegencyID           *int16    `gorm:"column:regency_id" json:"regency_id"`
	DistrictID          *int16    `gorm:"column:district_id" json:"district_id"`
	SubDistrictID       *int32    `gorm:"column:sub_district_id" json:"sub_district_id"`
	SrID                *int16    `gorm:"column:sr_id" json:"sr_id"`
	RayonID             *int16    `gorm:"column:rayon_id" json:"rayon_id"`
	BranchID            *int32    `gorm:"column:branch_id" json:"branch_id"`
	Latitude            *string   `gorm:"column:latitude" json:"latitude"`
	Longitude           *string   `gorm:"column:longitude" json:"longitude"`
}

// TableName Transaction's table name
func (*Transaction) TableName() string {
	return TableNameTransaction
}
