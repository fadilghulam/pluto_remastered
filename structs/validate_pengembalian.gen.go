package structs

const TableNameValidatePengembalian = "validate_pengembalian"

// ValidatePengembalian mapped from table <validate_pengembalian>
type ValidatePengembalian struct {
	ID                FlexibleString `gorm:"column:id;primaryKey;default:null" json:"id"`
	CustomerID        FlexibleString `gorm:"column:customer_id;not null" json:"customer_id"`
	PengembalianID    FlexibleString `gorm:"column:pengembalian_id;not null" json:"pengembalian_id"`
	UserID            int32          `gorm:"column:user_id;not null" json:"user_id"`
	Datetime          string         `gorm:"column:datetime;not null;default:now()" json:"datetime"`
	IsExist           int16          `gorm:"column:is_exist;not null" json:"is_exist"`
	ExistInfo         *string        `gorm:"column:exist_info;default:null" json:"exist_info"`
	IsValidQty        *int16         `gorm:"column:is_valid_qty;default:null" json:"is_valid_qty"`
	Qty               *string        `gorm:"column:qty;default:null" json:"qty"`
	IsValidDate       *int16         `gorm:"column:is_valid_date;default:null" json:"is_valid_date"`
	Date              *string        `gorm:"column:date;default:null" json:"date"`
	IsValidBrand      *int16         `gorm:"column:is_valid_brand;default:null" json:"is_valid_brand"`
	Brand             *string        `gorm:"column:brand;default:null" json:"brand"`
	CreatedAt         string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt         string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	LatitudeLongitude *string        `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
}

// TableName ValidatePengembalian's table name
func (*ValidatePengembalian) TableName() string {
	return TableNameValidatePengembalian
}
