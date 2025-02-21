package gpt_structs

const TableNameSubDistrict = "sub_district"

// SubDistrict mapped from table <sub_district>
type SubDistrict struct {
	ID         int64   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name       *string `gorm:"column:name" json:"name"`
	DistrictID *int32  `gorm:"column:district_id" json:"district_id"`
}

// TableName SubDistrict's table name
func (*SubDistrict) TableName() string {
	return TableNameSubDistrict
}
