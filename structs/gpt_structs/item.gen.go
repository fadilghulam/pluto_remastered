package gpt_structs

const TableNameItem = "item"

// Item mapped from table <item>
type Item struct {
	ID         int32   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name       *string `gorm:"column:name" json:"name"`
	ItemTypeID *int16  `gorm:"column:item_type_id" json:"item_type_id"`
}

// TableName Item's table name
func (*Item) TableName() string {
	return TableNameItem
}
