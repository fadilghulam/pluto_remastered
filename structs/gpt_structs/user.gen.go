package gpt_structs

import (
	"time"
)

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name       *string   `gorm:"column:name" json:"name"`
	UserTypeID *int16    `gorm:"column:user_type_id" json:"user_type_id"`
	IsActive   *int16    `gorm:"column:is_active" json:"is_active"`
	JoinAt     time.Time `gorm:"column:join_at" json:"join_at"`
	SrID       *int16    `gorm:"column:sr_id" json:"sr_id"`
	RayonID    *int16    `gorm:"column:rayon_id" json:"rayon_id"`
	BranchID   *int32    `gorm:"column:branch_id" json:"branch_id"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
