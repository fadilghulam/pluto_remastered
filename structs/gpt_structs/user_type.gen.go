// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameUserType = "user_type"

// UserType mapped from table <user_type>
type UserType struct {
	ID       int16  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Category string `gorm:"column:category" json:"category"`
}

// TableName UserType's table name
func (*UserType) TableName() string {
	return TableNameUserType
}
