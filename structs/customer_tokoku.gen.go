package structs

const TableNameCustomerTokoku = "customer_tokoku"

// CustomerTokoku mapped from table <customer_tokoku>
type CustomerTokoku struct {
	ID             FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID     FlexibleString `gorm:"column:customer_id;default:null" json:"customer_id"`
	UserID         *int64         `gorm:"column:user_id;default:null" json:"user_id"`
	Datetime       *string        `gorm:"column:datetime;default:null" json:"datetime"`
	CreatedAt      string         `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt      string         `gorm:"column:updated_at;default:now()" json:"updated_at"`
	SalesmanID     *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	TeamleaderID   *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	MerchandiserID *int32         `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	SyncKey        string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
}

// TableName CustomerTokoku's table name
func (*CustomerTokoku) TableName() string {
	return TableNameCustomerTokoku
}
