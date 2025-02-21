package structs

const TableNamePiutang = "piutang"

// Piutang mapped from table <piutang>
type Piutang struct {
	ID              FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	PenjualanID     FlexibleString `gorm:"column:penjualan_id;not null" json:"penjualan_id"`
	TanggalPiutang  string         `gorm:"column:tanggal_piutang;not null" json:"tanggal_piutang"`
	IsLunas         int32          `gorm:"column:is_lunas;not null" json:"is_lunas"`
	TanggalLunas    *string        `gorm:"column:tanggal_lunas;default:null" json:"tanggal_lunas"`
	DtmCrt          string         `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd          string         `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	TotalPiutang    float64        `gorm:"column:total_piutang;not null" json:"total_piutang"`
	CustomerID      FlexibleString `gorm:"column:customer_id;not null" json:"customer_id"`
	SyncKey         string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	IsLunas2        *int16         `gorm:"column:is_lunas2;default:null" json:"is_lunas2"`
	CustomerTipe    *string        `gorm:"column:customer_tipe;default:null" json:"customer_tipe"`
	CustomerTypeID  *int16         `gorm:"column:customer_type_id;default:null" json:"customer_type_id"`
	SalesmanTypeID  *int16         `gorm:"column:salesman_type_id;default:null" json:"salesman_type_id"`
	TeamleaderID    *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	UserID          *int32         `gorm:"column:user_id;default:null" json:"user_id"`
	UserIDSubtitute *int32         `gorm:"column:user_id_subtitute;default:null" json:"user_id_subtitute"`
}

// TableName Piutang's table name
func (*Piutang) TableName() string {
	return TableNamePiutang
}
