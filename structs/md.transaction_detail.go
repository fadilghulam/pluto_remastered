package structs

import (
	"time"
)

const TableNameMdTransactionDetail = "md.transaction_detail"

// Md Transaction Detail mapped from table <md.transaction_detail>
type MdTransactionDetail struct {
	ID            string    `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	TransactionID string    `gorm:"column:transaction_id;default:null" json:"transaction_id"`
	ItemID        int32     `gorm:"column:item_id;default:null" json:"item_id"`
	Qty           int64     `gorm:"column:qty;default:null" json:"qty"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	SyncKey       string    `gorm:"column:sync_key;default:now()" json:"sync_key"`
}

// TableName MdTransaction's table name
func (*MdTransactionDetail) TableName() string {
	return TableNameMdTransactionDetail
}
