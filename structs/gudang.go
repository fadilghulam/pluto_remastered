package structs

const TableNameGudang = "gudang"

// Gudang mapped from table <gudang>
type Gudang struct {
	ID          int16  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	SrID        int16  `gorm:"column:sr_id" json:"sr_id"`
	RayonID     int16  `gorm:"column:rayon_id" json:"rayon_id"`
	BranchID    int16  `gorm:"column:branch_id" json:"branch_id"`
	IsSalesman  int16  `gorm:"column:is_salesman; default: 1" json:"is_salesman"`
	Name        string `gorm:"column:name;not null" json:"name"`
	BranchIDOld int16  `gorm:"column:branch_id_old" json:"branch_id_old"`
	Deskripsi   string `gorm:"column:deskripsi" json:"deskripsi"`
}

// TableName Gudang's table name
func (*Gudang) TableName() string {
	return TableNameGudang
}
