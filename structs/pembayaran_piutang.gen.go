package structs

const TableNamePembayaranPiutang = "pembayaran_piutang"

// PembayaranPiutang mapped from table <pembayaran_piutang>
type PembayaranPiutang struct {
	ID                FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID        FlexibleString `gorm:"column:customer_id;not null" json:"customer_id"`
	SalesmanID        *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	TanggalPembayaran string         `gorm:"column:tanggal_pembayaran;not null" json:"tanggal_pembayaran"`
	DtmCrt            string         `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd            string         `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	TotalPembayaran   float64        `gorm:"column:total_pembayaran;not null" json:"total_pembayaran"`
	IsLunas           int32          `gorm:"column:is_lunas;not null;default:0" json:"is_lunas"`
	ImageNota         *string        `gorm:"column:image_nota;default:null" json:"image_nota"`
	TipePelunasan     int16          `gorm:"column:tipe_pelunasan;default:0" json:"tipe_pelunasan"`
	SyncKey           string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	IsComplete        *int16         `gorm:"column:is_complete;default:null" json:"is_complete"`
	PaymentID         FlexibleString `gorm:"column:payment_id;default:null" json:"payment_id"`
	BranchOldID       *int16         `gorm:"column:branch_old_id;default:null" json:"branch_old_id"`
	BranchID          *int16         `gorm:"column:branch_id;default:null" json:"branch_id"`
	SrID              *int16         `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID           *int16         `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	AreaID            *int32         `gorm:"column:area_id;default:null" json:"area_id"`
	SalesmanTipe      *string        `gorm:"column:salesman_tipe;default:null" json:"salesman_tipe"`
	LatitudeLongitude *string        `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
	PengembalianID    FlexibleString `gorm:"column:pengembalian_id;default:null" json:"pengembalian_id"`
	CustomerTipe      *string        `gorm:"column:customer_tipe;default:null" json:"customer_tipe"`
	CustomerTypeID    *int16         `gorm:"column:customer_type_id;default:null" json:"customer_type_id"`
	SalesmanTypeID    *int16         `gorm:"column:salesman_type_id;default:null" json:"salesman_type_id"`
	MerchandiserID    *int32         `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	ToRefID           *int32         `gorm:"column:to_ref_id;default:null" json:"to_ref_id"`
	ToRefName         *string        `gorm:"column:to_ref_name;default:null" json:"to_ref_name"`
	TeamleaderID      *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	UserID            *int32         `gorm:"column:user_id;default:null" json:"user_id"`
	UserIDSubtitute   *int32         `gorm:"column:user_id_subtitute;default:null" json:"user_id_subtitute"`
}

// TableName PembayaranPiutang's table name
func (*PembayaranPiutang) TableName() string {
	return TableNamePembayaranPiutang
}
