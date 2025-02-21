package structs

const TableNameValidateKunjungan = "validate_kunjungan"

// ValidateKunjungan mapped from table <validate_kunjungan>
type ValidateKunjungan struct {
	ID                FlexibleString `gorm:"column:id;primaryKey;default:null" json:"id"`
	CustomerID        FlexibleString `gorm:"column:customer_id;not null" json:"customer_id"`
	KunjunganID       FlexibleString `gorm:"column:kunjungan_id;not null" json:"kunjungan_id"`
	UserID            int32          `gorm:"column:user_id;not null" json:"user_id"`
	Datetime          string         `gorm:"column:datetime;not null;default:now()" json:"datetime"`
	IsExist           int16          `gorm:"column:is_exist;not null" json:"is_exist"`
	ExistInfo         *string        `gorm:"column:exist_info;default:null" json:"exist_info"`
	IsValid           *int16         `gorm:"column:is_valid;default:null" json:"is_valid"`
	ValidateInfo      *string        `gorm:"column:validate_info;default:null" json:"validate_info"`
	CreatedAt         string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt         string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	LatitudeLongitude *string        `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
}

// TableName ValidateKunjungan's table name
func (*ValidateKunjungan) TableName() string {
	return TableNameValidateKunjungan
}
