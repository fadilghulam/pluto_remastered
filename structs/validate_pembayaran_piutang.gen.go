package structs

const TableNameValidatePembayaranPiutang = "validate_pembayaran_piutang"

// ValidatePembayaranPiutang mapped from table <validate_pembayaran_piutang>
type ValidatePembayaranPiutang struct {
	ID                   FlexibleString `gorm:"column:id;primaryKey;default:null" json:"id"`
	CustomerID           FlexibleString `gorm:"column:customer_id;not null" json:"customer_id"`
	PembayaranPiutangID  FlexibleString `gorm:"column:pembayaran_piutang_id;not null" json:"pembayaran_piutang_id"`
	UserID               int32          `gorm:"column:user_id;not null" json:"user_id"`
	IsExist              int16          `gorm:"column:is_exist;not null" json:"is_exist"`
	ExistInfo            *string        `gorm:"column:exist_info;default:null" json:"exist_info"`
	IsValidNominal       *int16         `gorm:"column:is_valid_nominal;default:null" json:"is_valid_nominal"`
	Nominal              *int32         `gorm:"column:nominal;default:null" json:"nominal"`
	IsValidDate          *int16         `gorm:"column:is_valid_date;default:null" json:"is_valid_date"`
	Date                 *string        `gorm:"column:date;default:null" json:"date"`
	IsValidPaymentMethod *int16         `gorm:"column:is_valid_payment_method;default:null" json:"is_valid_payment_method"`
	PaymentMethod        *string        `gorm:"column:payment_method;default:null" json:"payment_method"`
	CreatedAt            string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt            string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	Datetime             string         `gorm:"column:datetime;not null;default:now()" json:"datetime"`
	LatitudeLongitude    *string        `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
}

// TableName ValidatePembayaranPiutang's table name
func (*ValidatePembayaranPiutang) TableName() string {
	return TableNameValidatePembayaranPiutang
}
