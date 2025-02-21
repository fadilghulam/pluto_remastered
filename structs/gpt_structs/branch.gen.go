package gpt_structs

const TableNameBranch = "branch"

// Branch mapped from table <branch>
type Branch struct {
	ID      int32   `gorm:"column:id;primaryKey;autoIncrement:true;comment:this is for branch office identifier / id" json:"id"` // this is for branch office identifier / id
	Name    *string `gorm:"column:name" json:"name"`
	RayonID *int16  `gorm:"column:rayon_id" json:"rayon_id"`
}

// TableName Branch's table name
func (*Branch) TableName() string {
	return TableNameBranch
}
