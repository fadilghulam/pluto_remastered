package structs

const TableNameValidateCustomer = "validate_customer"

// ValidateCustomer mapped from table <validate_customer>
type ValidateCustomer struct {
	ID                       FlexibleString  `gorm:"column:id;primaryKey;default:null" json:"id"`
	CustomerID               FlexibleString  `gorm:"column:customer_id;not null" json:"customer_id"`
	UserID                   int32           `gorm:"column:user_id;not null" json:"user_id"`
	Datetime                 string          `gorm:"column:datetime;not null;default:now()" json:"datetime"`
	IsExist                  int16           `gorm:"column:is_exist;not null" json:"is_exist"`
	ExistInfo                *string         `gorm:"column:exist_info;default:null" json:"exist_info"`
	IsDuplicate              *int16          `gorm:"column:is_duplicate;default:null" json:"is_duplicate"`
	CustomerIdsDuplicate     *StringArray    `gorm:"column:customer_ids_duplicate;default:null" json:"customer_ids_duplicate"`
	CustomerIDDuplicateValid *FlexibleString `gorm:"column:customer_id_duplicate_valid;default:null" json:"customer_id_duplicate_valid"`
	IsValidCustomerName      *int16          `gorm:"column:is_valid_customer_name;default:null" json:"is_valid_customer_name"`
	CustomerName             *string         `gorm:"column:customer_name;default:null" json:"customer_name"`
	IsValidOutletName        *int16          `gorm:"column:is_valid_outlet_name;default:null" json:"is_valid_outlet_name"`
	OutletName               *string         `gorm:"column:outlet_name;default:null" json:"outlet_name"`
	IsValidPhone             *int16          `gorm:"column:is_valid_phone;default:null" json:"is_valid_phone"`
	Phone                    *string         `gorm:"column:phone;default:null" json:"phone"`
	IsValidType              *int16          `gorm:"column:is_valid_type;default:null" json:"is_valid_type"`
	Type                     *int32          `gorm:"column:type;default:null" json:"type"`
	IsValidLatlong           *int16          `gorm:"column:is_valid_latlong;default:null" json:"is_valid_latlong"`
	Latlong                  *string         `gorm:"column:latlong;default:null" json:"latlong"`
	IsValidAddress           *int16          `gorm:"column:is_valid_address;default:null" json:"is_valid_address"`
	Alamat                   *string         `gorm:"column:alamat;default:null" json:"alamat"`
	KelurahanID              *int64          `gorm:"column:kelurahan_id;default:null" json:"kelurahan_id"`
	KelurahanName            *string         `gorm:"column:kelurahan_name;default:null" json:"kelurahan_name"`
	KecamatanID              *int32          `gorm:"column:kecamatan_id;default:null" json:"kecamatan_id"`
	KecamatanName            *string         `gorm:"column:kecamatan_name;default:null" json:"kecamatan_name"`
	KabupatenID              *int32          `gorm:"column:kabupaten_id;default:null" json:"kabupaten_id"`
	KabupatenName            *string         `gorm:"column:kabupaten_name;default:null" json:"kabupaten_name"`
	ProvinsiID               *int32          `gorm:"column:provinsi_id;default:null" json:"provinsi_id"`
	ProvinsiName             *string         `gorm:"column:provinsi_name;default:null" json:"provinsi_name"`
	CreatedAt                string          `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt                string          `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	LatitudeLongitude        *string         `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
}

// TableName ValidateCustomer's table name
func (*ValidateCustomer) TableName() string {
	return TableNameValidateCustomer
}
