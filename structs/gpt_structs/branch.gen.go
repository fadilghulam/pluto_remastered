// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameBranch = "branch"

// Branch mapped from table <branch>
type Branch struct {
	ID   int32  `gorm:"column:id;primaryKey;autoIncrement:true;comment:this is for branch office identifier / id" json:"id"` // this is for branch office identifier / id
	Name string `gorm:"column:name" json:"name"`
}

// TableName Branch's table name
func (*Branch) TableName() string {
	return TableNameBranch
}
