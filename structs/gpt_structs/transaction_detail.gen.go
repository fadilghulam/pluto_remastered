package gpt_structs

const TableNameTransactionDetail = "transaction_detail"

// TransactionDetail mapped from table <transaction_detail>
type TransactionDetail struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	TransactionID *int64 `gorm:"column:transaction_id" json:"transaction_id"`
	ItemID        *int32 `gorm:"column:item_id" json:"item_id"`
	Qty           *int16 `gorm:"column:qty" json:"qty"`
	Price         *int64 `gorm:"column:price" json:"price"`
	SubTotal      *int64 `gorm:"column:sub_total" json:"sub_total"`
}

// TableName TransactionDetail's table name
func (*TransactionDetail) TableName() string {
	return TableNameTransactionDetail
}
