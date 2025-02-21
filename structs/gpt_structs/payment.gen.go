package gpt_structs

import (
	"time"
)

const TableNamePayment = "payment"

// Payment mapped from table <payment>
type Payment struct {
	ID                int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	TransactionID     *int64    `gorm:"column:transaction_id" json:"transaction_id"`
	InputAt           time.Time `gorm:"column:input_at" json:"input_at"`
	UserIDExecutor    *int32    `gorm:"column:user_id_executor" json:"user_id_executor"`
	UserIDHolder      *int32    `gorm:"column:user_id_holder" json:"user_id_holder"`
	CustomerID        *int64    `gorm:"column:customer_id" json:"customer_id"`
	Amount            *int64    `gorm:"column:amount" json:"amount"`
	PaymentMethod     *string   `gorm:"column:payment_method" json:"payment_method"`
	BankName          *string   `gorm:"column:bank_name" json:"bank_name"`
	BankAccountNumber *string   `gorm:"column:bank_account_number" json:"bank_account_number"`
	BankAccountName   *string   `gorm:"column:bank_account_name" json:"bank_account_name"`
	IsPaid            *int16    `gorm:"column:is_paid" json:"is_paid"`
	PaidAt            time.Time `gorm:"column:paid_at" json:"paid_at"`
	SrID              *int16    `gorm:"column:sr_id" json:"sr_id"`
	RayonID           *int16    `gorm:"column:rayon_id" json:"rayon_id"`
	BranchID          *int32    `gorm:"column:branch_id" json:"branch_id"`
	Latitude          *string   `gorm:"column:latitude" json:"latitude"`
	Longitude         *string   `gorm:"column:longitude" json:"longitude"`
}

// TableName Payment's table name
func (*Payment) TableName() string {
	return TableNamePayment
}
