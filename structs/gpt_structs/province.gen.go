package gpt_structs

const TableNameProvince = "province"

// Province mapped from table <province>
type Province struct {
	ID   int16   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name *string `gorm:"column:name" json:"name"`
}

// TableName Province's table name
func (*Province) TableName() string {
	return TableNameProvince
}
