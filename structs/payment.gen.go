package structs

const TableNamePayment = "payment"

// Payment mapped from table <payment>
type Payment struct {
	ID                FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID        FlexibleString `gorm:"column:customer_id;not null" json:"customer_id"`
	SalesmanID        *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	PenjualanID       FlexibleString `gorm:"column:penjualan_id;not null" json:"penjualan_id"`
	NamaPenyetor      string         `gorm:"column:nama_penyetor;not null" json:"nama_penyetor"`
	NoReferensi       *string        `gorm:"column:no_referensi;default:null" json:"no_referensi"`
	Tipe              string         `gorm:"column:tipe;not null" json:"tipe"`
	Nominal           int64          `gorm:"column:nominal;not null" json:"nominal"`
	TanggalJatuhTempo string         `gorm:"column:tanggal_jatuh_tempo;not null;default:now()" json:"tanggal_jatuh_tempo"`
	IsCair            int16          `gorm:"column:is_cair;not null;default:1" json:"is_cair"`
	TanggalCair       *string        `gorm:"column:tanggal_cair;default:null" json:"tanggal_cair"`
	VerifBy           *int32         `gorm:"column:verif_by;default:null" json:"verif_by"`
	TanggalVerif      *string        `gorm:"column:tanggal_verif;default:null" json:"tanggal_verif"`
	DtmCrt            string         `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd            string         `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	TanggalTransaksi  *string        `gorm:"column:tanggal_transaksi;default:null" json:"tanggal_transaksi"`
	Bank              *string        `gorm:"column:bank;default:null" json:"bank"`
	Keterangan        *string        `gorm:"column:keterangan;default:null" json:"keterangan"`
	SyncKey           string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	BuktiBayar        *string        `gorm:"column:bukti_bayar;default:null" json:"bukti_bayar"`
	BranchIDOld       *int16         `gorm:"column:branch_id_old;default:null" json:"branch_id_old"`
	IsVerif           int16          `gorm:"column:is_verif;default:0" json:"is_verif"`
	BranchID          *int16         `gorm:"column:branch_id;default:null" json:"branch_id"`
	SrID              *int16         `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID           *int16         `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	AreaID            *int32         `gorm:"column:area_id;default:null" json:"area_id"`
	SalesmanTipe      *string        `gorm:"column:salesman_tipe;default:null" json:"salesman_tipe"`
	AccountIDTujuan   *int16         `gorm:"column:account_id_tujuan;default:null" json:"account_id_tujuan"`
	PengembalianID    FlexibleString `gorm:"column:pengembalian_id;default:null" json:"pengembalian_id"`
	CustomerTipe      *string        `gorm:"column:customer_tipe;default:null" json:"customer_tipe"`
	Provinsi          *string        `gorm:"column:provinsi;default:null" json:"provinsi"`
	Kabupaten         *string        `gorm:"column:kabupaten;default:null" json:"kabupaten"`
	Kecamatan         *string        `gorm:"column:kecamatan;default:null" json:"kecamatan"`
	Kelurahan         *string        `gorm:"column:kelurahan;default:null" json:"kelurahan"`
	CustomerTypeID    *int16         `gorm:"column:customer_type_id;default:null" json:"customer_type_id"`
	SalesmanTypeID    *int16         `gorm:"column:salesman_type_id;default:null" json:"salesman_type_id"`
	MerchandiserID    *int32         `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	ToRefID           *int32         `gorm:"column:to_ref_id;default:null" json:"to_ref_id"`
	ToRefName         *string        `gorm:"column:to_ref_name;default:null" json:"to_ref_name"`
	TeamleaderID      *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	UserID            *int32         `gorm:"column:user_id;default:null" json:"user_id"`
	UserIDSubtitute   *int32         `gorm:"column:user_id_subtitute;default:null" json:"user_id_subtitute"`
}

// TableName Payment's table name
func (*Payment) TableName() string {
	return TableNamePayment
}
