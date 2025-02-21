package gpt_structs

const TableNameRegency = "regency"

// Regency mapped from table <regency>
type Regency struct {
	ID         int16   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name       *string `gorm:"column:name" json:"name"`
	ProvinceID *int16  `gorm:"column:province_id" json:"province_id"`
}

// TableName Regency's table name
func (*Regency) TableName() string {
	return TableNameRegency
}
