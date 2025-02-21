package gpt_structs

const TableNameTransactionType = "transaction_type"

// TransactionType mapped from table <transaction_type>
type TransactionType struct {
	ID       int16   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name     *string `gorm:"column:name" json:"name"`
	Category *string `gorm:"column:category" json:"category"`
}

// TableName TransactionType's table name
func (*TransactionType) TableName() string {
	return TableNameTransactionType
}
