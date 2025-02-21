package gpt_structs

const TableNameSr = "sr"

// Sr mapped from table <sr>
type Sr struct {
	ID   int16   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name *string `gorm:"column:name" json:"name"`
}

// TableName Sr's table name
func (*Sr) TableName() string {
	return TableNameSr
}
