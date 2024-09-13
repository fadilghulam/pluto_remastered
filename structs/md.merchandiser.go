package structs

import "time"

const TableNameMerchandiser = "md.merchandiser"

// Merchandiser mapped from table <md.merchandiser>
type Merchandiser struct {
	ID           int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID       int32      `gorm:"column:user_id" json:"user_id"`
	Name         string     `gorm:"column:name" json:"name"`
	SrID         int16      `gorm:"column:sr_id" json:"sr_id"`
	RayonID      int16      `gorm:"column:rayon_id" json:"rayon_id"`
	BranchID     int16      `gorm:"column:branch_id" json:"branch_id"`
	AreaID       Int32Array `gorm:"column:area_id" json:"area_id"`
	CreatedAt    time.Time  `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;default:now()" json:"updated_at"`
	PromotorId   int16      `gorm:"column:promotor_id;default: null" json:"promotor_id"`
	TeamleaderId int16      `gorm:"column:teamleader_id;default: null" json:"teamleader_id"`
	IsAktif      int16      `gorm:"column:is_aktif;default: 1" json:"is_aktif"`
	SalesmanId   int32      `gorm:"column:salesman_id;default: null" json:"salesman_id"`
	SkID         string     `gorm:"column:sk_id" json:"sk_id"`
	TipeSalesman string     `gorm:"column:tipe_salesman;default: null" json:"tipe_salesman"`
}

// TableName Salesman's table name
func (*Merchandiser) TableName() string {
	return TableNameMerchandiser
}
