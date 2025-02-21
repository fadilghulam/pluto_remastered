package structs

const TableNameValidatePenjualan = "validate_penjualan"

// ValidatePenjualan mapped from table <validate_penjualan>
type ValidatePenjualan struct {
	ID                   FlexibleString `gorm:"column:id;primaryKey;default:null" json:"id"`
	CustomerID           FlexibleString `gorm:"column:customer_id;not null" json:"customer_id"`
	PenjualanID          FlexibleString `gorm:"column:penjualan_id;not null" json:"penjualan_id"`
	UserID               int32          `gorm:"column:user_id;not null" json:"user_id"`
	Datetime             string         `gorm:"column:datetime;not null;default:now()" json:"datetime"`
	IsExist              int16          `gorm:"column:is_exist;not null" json:"is_exist"`
	ExistInfo            *string        `gorm:"column:exist_info;default:null" json:"exist_info"`
	IsValidType          *int16         `gorm:"column:is_valid_type;default:null" json:"is_valid_type"`
	Type                 *string        `gorm:"column:type;default:null" json:"type"`
	IsValidPrice         *int16         `gorm:"column:is_valid_price;default:null" json:"is_valid_price"`
	Price                *string        `gorm:"column:price;default:null" json:"price"`
	IsValidQty           *int16         `gorm:"column:is_valid_qty;default:null" json:"is_valid_qty"`
	Qty                  *string        `gorm:"column:qty;default:null" json:"qty"`
	IsValidDate          *int16         `gorm:"column:is_valid_date;default:null" json:"is_valid_date"`
	Date                 *string        `gorm:"column:date;default:null" json:"date"`
	IsValidPaymentMethod *int16         `gorm:"column:is_valid_payment_method;default:null" json:"is_valid_payment_method"`
	PaymentMethod        *string        `gorm:"column:payment_method;default:null" json:"payment_method"`
	IsValidBrand         *int16         `gorm:"column:is_valid_brand;default:null" json:"is_valid_brand"`
	Brand                *string        `gorm:"column:brand;default:null" json:"brand"`
	IsValidProgram       *int16         `gorm:"column:is_valid_program;default:null" json:"is_valid_program"`
	Program              *string        `gorm:"column:program;default:null" json:"program"`
	CreatedAt            string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt            string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	LatitudeLongitude    *string        `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
}

// TableName ValidatePenjualan's table name
func (*ValidatePenjualan) TableName() string {
	return TableNameValidatePenjualan
}
