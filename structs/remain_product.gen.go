// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package structs

import (
	"time"
)

const TableNameRemainProduct = "remain_product"

// RemainProduct mapped from table <remain_product>
type RemainProduct struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	SalesmanID int32     `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	CustomerID int64     `gorm:"column:customer_id;default:null" json:"customer_id"`
	Datetime   time.Time `gorm:"column:datetime;default:null" json:"datetime"`
	Photos     string    `gorm:"column:photos;default:null" json:"photos"`
	Note       string    `gorm:"column:note;default:null" json:"note"`
	SrID       int16     `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID    int16     `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	BranchID   int16     `gorm:"column:branch_id;default:null" json:"branch_id"`
	AreaID     int32     `gorm:"column:area_id;default:null" json:"area_id"`
	CreatedAt  time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:now()" json:"updated_at"`
	SyncKey    string    `gorm:"column:sync_key;default:now()" json:"sync_key"`
}

// TableName RemainProduct's table name
func (*RemainProduct) TableName() string {
	return TableNameRemainProduct
}
