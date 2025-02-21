package gpt_structs

const TableNameCustomerType = "customer_type"

// CustomerType mapped from table <customer_type>
type CustomerType struct {
	ID       int16   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name     *string `gorm:"column:name" json:"name"`
	Category *string `gorm:"column:category" json:"category"`
}

// TableName CustomerType's table name
func (*CustomerType) TableName() string {
	return TableNameCustomerType
}
