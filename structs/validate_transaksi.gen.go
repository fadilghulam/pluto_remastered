package structs

const TableNameValidateTransaksi = "validate_transaksi"

// ValidateTransaksi mapped from table <validate_transaksi>
type ValidateTransaksi struct {
	ID                FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID        FlexibleString `gorm:"column:customer_id;default:null" json:"customer_id"`
	UserID            *int32         `gorm:"column:user_id;default:null" json:"user_id"`
	UserLevelID       *int32         `gorm:"column:user_level_id;default:null" json:"user_level_id"`
	SrID              *int16         `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID           *int16         `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	BranchID          *int16         `gorm:"column:branch_id;default:null" json:"branch_id"`
	TransactionID     FlexibleString `gorm:"column:transaction_id;default:null" json:"transaction_id"`
	TransactionType   FlexibleString `gorm:"column:transaction_type;default:null" json:"transaction_type"`
	Datetime          *string        `gorm:"column:datetime;default:null" json:"datetime"`
	IsValid           *int16         `gorm:"column:is_valid;default:null" json:"is_valid"`
	Note              *string        `gorm:"column:note;default:null" json:"note"`
	LatitudeLongitude *string        `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
	CreatedAt         string         `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt         string         `gorm:"column:updated_at;default:now()" json:"updated_at"`
}

// TableName ValidateTransaksi's table name
func (*ValidateTransaksi) TableName() string {
	return TableNameValidateTransaksi
}
