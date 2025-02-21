package structs

const TableNameRemainProductDetail = "remain_product_detail"

// RemainProductDetail mapped from table <remain_product_detail>
type RemainProductDetail struct {
	ID              FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	RemainProductID FlexibleString `gorm:"column:remain_product_id;default:null" json:"remain_product_id"`
	ProductID       *int32         `gorm:"column:product_id;default:null" json:"product_id"`
	Qty             *int32         `gorm:"column:qty;default:null" json:"qty"`
	CreatedAt       string         `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt       string         `gorm:"column:updated_at;default:now()" json:"updated_at"`
	SyncKey         string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
}

// TableName RemainProductDetail's table name
func (*RemainProductDetail) TableName() string {
	return TableNameRemainProductDetail
}
