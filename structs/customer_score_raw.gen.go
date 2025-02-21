package structs

const TableNameCustomerScoreRaw = "customer_score_raw"

// CustomerScoreRaw mapped from table <customer_score_raw>
type CustomerScoreRaw struct {
	ID          FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID  FlexibleString `gorm:"column:customer_id;default:null" json:"customer_id"`
	SalesmanID  *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	BranchID    *int32         `gorm:"column:branch_id;default:null" json:"branch_id"`
	Score       *int16         `gorm:"column:score;default:null" json:"score"`
	Indicator   *int16         `gorm:"column:indicator;default:null" json:"indicator"`
	CreatedAt   string         `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt   string         `gorm:"column:updated_at;default:now()" json:"updated_at"`
	DateJoin    *string        `gorm:"column:date_join;default:null" json:"date_join"`
	BySystem    *int16         `gorm:"column:by_system;default:null" json:"by_system"`
	DateCreated string         `gorm:"column:date_created;default:CURRENT_DATE" json:"date_created"`
}

// TableName CustomerScoreRaw's table name
func (*CustomerScoreRaw) TableName() string {
	return TableNameCustomerScoreRaw
}
