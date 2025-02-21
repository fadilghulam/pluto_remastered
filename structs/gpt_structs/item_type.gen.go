package gpt_structs

const TableNameItemType = "item_type"

// ItemType mapped from table <item_type>
type ItemType struct {
	ID   int16   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name *string `gorm:"column:name" json:"name"`
}

// TableName ItemType's table name
func (*ItemType) TableName() string {
	return TableNameItemType
}
