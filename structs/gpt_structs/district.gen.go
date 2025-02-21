package gpt_structs

const TableNameDistrict = "district"

// District mapped from table <district>
type District struct {
	ID        int32   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name      *string `gorm:"column:name" json:"name"`
	RegencyID *int16  `gorm:"column:regency_id" json:"regency_id"`
}

// TableName District's table name
func (*District) TableName() string {
	return TableNameDistrict
}
