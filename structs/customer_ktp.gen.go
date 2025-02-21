package structs

const TableNameCustomerKtp = "customer_ktp"

// CustomerKtp mapped from table <customer_ktp>
type CustomerKtp struct {
	ID             FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID     FlexibleString `gorm:"column:customer_id;default:null" json:"customer_id"`
	ImageKtp       *string        `gorm:"column:image_ktp;default:null" json:"image_ktp"`
	CreatedAt      string         `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt      string         `gorm:"column:updated_at;default:now()" json:"updated_at"`
	SalesmanID     *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	TeamleaderID   *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	MerchandiserID *int32         `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	SyncKey        string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	NameKtp        *string        `gorm:"column:name_ktp;default:null" json:"name_ktp"`
	NikKtp         *string        `gorm:"column:nik_ktp;default:null" json:"nik_ktp"`
}

// TableName CustomerKtp's table name
func (*CustomerKtp) TableName() string {
	return TableNameCustomerKtp
}
