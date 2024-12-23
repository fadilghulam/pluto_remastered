package gpt_structs

const TableNameDictionary = "dictionary"

// Dictionary mapped from table <dictionary>
type Dictionary struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name            string `gorm:"column:name" json:"name"`
	Url             string `gorm:"column:url" json:"url"`
	Params          string `gorm:"column:params" json:"params"`
	Tags            string `gorm:"column:tags" json:"tags"`
	Description     string `gorm:"column:description" json:"description"`
	DataInformation string `gorm:"column:data_information" json:"data_information"`
	Method          string `gorm:"column:method" json:"method"`
}

// TableName Dictionary's table name
func (*Dictionary) TableName() string {
	return TableNameDictionary
}
