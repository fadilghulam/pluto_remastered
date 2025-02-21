package structs

const TableNamePenjualan = "penjualan"

// Penjualan mapped from table <penjualan>
type Penjualan struct {
	ID                FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	SalesmanID        *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	CustomerID        FlexibleString `gorm:"column:customer_id;default:null" json:"customer_id"`
	IsKredit          int32          `gorm:"column:is_kredit;not null" json:"is_kredit"`
	TipePenjualan     int32          `gorm:"column:tipe_penjualan;not null" json:"tipe_penjualan"`
	TanggalPenjualan  string         `gorm:"column:tanggal_penjualan;not null" json:"tanggal_penjualan"`
	DtmCrt            string         `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd            string         `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	TotalPenjualan    float64        `gorm:"column:total_penjualan;not null" json:"total_penjualan"`
	ImageNota         *string        `gorm:"column:image_nota;default:null" json:"image_nota"`
	NoNota            *string        `gorm:"column:no_nota;default:null" json:"no_nota"`
	SyncKey           string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	BranchIDOld       *int16         `gorm:"column:branch_id_old;default:null" json:"branch_id_old"`
	BranchID          *int16         `gorm:"column:branch_id;default:null" json:"branch_id"`
	SrID              *int16         `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID           *int16         `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	AreaID            *int32         `gorm:"column:area_id;default:null" json:"area_id"`
	SalesmanTipe      *string        `gorm:"column:salesman_tipe;default:null" json:"salesman_tipe"`
	LatitudeLongitude *string        `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
	TeamleaderID      *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	MerchandiserID    *int32         `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	CustomerTipe      *string        `gorm:"column:customer_tipe;default:null" json:"customer_tipe"`
	Provinsi          *string        `gorm:"column:provinsi;default:null" json:"provinsi"`
	Kabupaten         *string        `gorm:"column:kabupaten;default:null" json:"kabupaten"`
	Kecamatan         *string        `gorm:"column:kecamatan;default:null" json:"kecamatan"`
	Kelurahan         *string        `gorm:"column:kelurahan;default:null" json:"kelurahan"`
	CustomerTypeID    *int16         `gorm:"column:customer_type_id;default:null" json:"customer_type_id"`
	SalesmanTypeID    *int16         `gorm:"column:salesman_type_id;default:null" json:"salesman_type_id"`
	ImageNotaPrint    *string        `gorm:"column:image_nota_print;default:null" json:"image_nota_print"`
	ToRefID           *int32         `gorm:"column:to_ref_id;default:null" json:"to_ref_id"`
	ToRefName         *string        `gorm:"column:to_ref_name;default:null" json:"to_ref_name"`
	ImageBuktiSerah   *string        `gorm:"column:image_bukti_serah;default:null" json:"image_bukti_serah"`
	OutletID          FlexibleString `gorm:"column:outlet_id;default:null" json:"outlet_id"`
	IsDoubleCredit    int16          `gorm:"column:is_double_credit;not null;default:0" json:"is_double_credit"`
	PlafonDiff        float64        `gorm:"column:plafon_diff;not null;default:0" json:"plafon_diff"`
	UserID            *int32         `gorm:"column:user_id;default:null" json:"user_id"`
	UserIDSubtitute   *int32         `gorm:"column:user_id_subtitute;default:null" json:"user_id_subtitute"`
}

// TableName Penjualan's table name
func (*Penjualan) TableName() string {
	return TableNamePenjualan
}
