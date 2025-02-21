package gpt_structs

const TableNameRayon = "rayon"

// Rayon mapped from table <rayon>
type Rayon struct {
	ID   int16   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name *string `gorm:"column:name" json:"name"`
	SrID *int16  `gorm:"column:sr_id" json:"sr_id"`
}

// TableName Rayon's table name
func (*Rayon) TableName() string {
	return TableNameRayon
}
