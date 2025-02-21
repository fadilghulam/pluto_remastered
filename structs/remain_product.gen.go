package structs

const TableNameRemainProduct = "remain_product"

// RemainProduct mapped from table <remain_product>
type RemainProduct struct {
	ID         FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	SalesmanID *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	CustomerID FlexibleString `gorm:"column:customer_id;default:null" json:"customer_id"`
	Datetime   *string        `gorm:"column:datetime;default:null" json:"datetime"`
	Photos     *string        `gorm:"column:photos;default:null" json:"photos"`
	Note       *string        `gorm:"column:note;default:null" json:"note"`
	SrID       *int16         `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID    *int16         `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	BranchID   *int16         `gorm:"column:branch_id;default:null" json:"branch_id"`
	AreaID     *int32         `gorm:"column:area_id;default:null" json:"area_id"`
	CreatedAt  string         `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt  string         `gorm:"column:updated_at;default:now()" json:"updated_at"`
	SyncKey    string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
}

// TableName RemainProduct's table name
func (*RemainProduct) TableName() string {
	return TableNameRemainProduct
}
