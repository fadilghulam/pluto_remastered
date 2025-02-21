package structs

const TableNameQrCodeHistory = "qr_code_history"

// QrCodeHistory mapped from table <qr_code_history>
type QrCodeHistory struct {
	ID             FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID     FlexibleString `gorm:"column:customer_id;not null" json:"customer_id"`
	UserID         int32          `gorm:"column:user_id;not null" json:"user_id"`
	Action         string         `gorm:"column:action;not null;comment:PEMASANGAN BARU, PENGGANTIAN" json:"action"` // PEMASANGAN BARU, PENGGANTIAN
	Note           *string        `gorm:"column:note;default:null" json:"note"`
	PhotoBefore    *string        `gorm:"column:photo_before;default:null" json:"photo_before"`
	PhotoAfter     *string        `gorm:"column:photo_after;default:null" json:"photo_after"`
	CreatedAt      string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt      string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	Datetime       string         `gorm:"column:datetime;not null" json:"datetime"`
	LatLong        *string        `gorm:"column:lat_long;default:null" json:"lat_long"`
	QrCode         string         `gorm:"column:qr_code;not null" json:"qr_code"`
	IsMockLocation int16          `gorm:"column:is_mock_location;not null;default:0" json:"is_mock_location"`
	SyncKey        string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	QrCodeExisting *string        `gorm:"column:qr_code_existing;default:null" json:"qr_code_existing"`
	SalesmanID     *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	CustomerTipe   *int16         `gorm:"column:customer_tipe;default:null" json:"customer_tipe"`
	TeamleaderID   *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
}

// TableName QrCodeHistory's table name
func (*QrCodeHistory) TableName() string {
	return TableNameQrCodeHistory
}
