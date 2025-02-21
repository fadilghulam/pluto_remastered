package structs

const TableNameKunjungan = "kunjungan"

// Kunjungan mapped from table <kunjungan>
type Kunjungan struct {
	ID                    FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	SalesmanID            *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	CustomerID            FlexibleString `gorm:"column:customer_id;default:null" json:"customer_id"`
	TanggalKunjungan      string         `gorm:"column:tanggal_kunjungan;not null" json:"tanggal_kunjungan"`
	DtmCrt                string         `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd                string         `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	StatusToko            *string        `gorm:"column:status_toko;default:null" json:"status_toko"`
	Keterangan            *string        `gorm:"column:keterangan;default:null" json:"keterangan"`
	SyncKey               string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	LatitudeLongitude     *string        `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
	ImageKunjungan        *string        `gorm:"column:image_kunjungan;default:null" json:"image_kunjungan"`
	BranchIDOld           *int16         `gorm:"column:branch_id_old;default:null" json:"branch_id_old"`
	IsMockLocation        *int16         `gorm:"column:is_mock_location;default:null" json:"is_mock_location"`
	TanggalImage          *string        `gorm:"column:tanggal_image;default:null" json:"tanggal_image"`
	IsRoute               int16          `gorm:"column:is_route;not null;default:0" json:"is_route"`
	AmbilGrosirID         string         `gorm:"column:ambil_grosir_id;default:-1;comment:-1/null tidak ambil di grosir0 ambil di grosir yg belum terdaftar>0 ambil di grosir yg sudah terdaftar" json:"ambil_grosir_id"` // -1/null tidak ambil di grosir0 ambil di grosir yg belum terdaftar>0 ambil di grosir yg sudah terdaftar
	AmbilGrosirKeterangan *string        `gorm:"column:ambil_grosir_keterangan;default:null" json:"ambil_grosir_keterangan"`
	BranchID              *int16         `gorm:"column:branch_id;default:null" json:"branch_id"`
	SrID                  *int16         `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID               *int16         `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	AreaID                *int32         `gorm:"column:area_id;default:null" json:"area_id"`
	SalesmanTipe          *string        `gorm:"column:salesman_tipe;default:null" json:"salesman_tipe"`
	CustomerTipe          *string        `gorm:"column:customer_tipe;default:null" json:"customer_tipe"`
	AllowCheckin          *int16         `gorm:"column:allow_checkin;default:null" json:"allow_checkin"`
	CustomerTypeID        *int16         `gorm:"column:customer_type_id;default:null" json:"customer_type_id"`
	SalesmanTypeID        *int16         `gorm:"column:salesman_type_id;default:null" json:"salesman_type_id"`
	MerchandiserID        *int32         `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	TeamleaderID          *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	SpvID                 *int32         `gorm:"column:spv_id;default:null" json:"spv_id"`
	Provinsi              *string        `gorm:"column:provinsi;default:null" json:"provinsi"`
	Kabupaten             *string        `gorm:"column:kabupaten;default:null" json:"kabupaten"`
	Kecamatan             *string        `gorm:"column:kecamatan;default:null" json:"kecamatan"`
	Kelurahan             *string        `gorm:"column:kelurahan;default:null" json:"kelurahan"`
	CustomerName          *string        `gorm:"column:customer_name;default:null" json:"customer_name"`
	OutletID              FlexibleString `gorm:"column:outlet_id;default:null" json:"outlet_id"`
	SubjectTypeID         *int32         `gorm:"column:subject_type_id;default:null" json:"subject_type_id"`
	OutletTypeID          *int32         `gorm:"column:outlet_type_id;default:null" json:"outlet_type_id"`
	ToRefID               *int32         `gorm:"column:to_ref_id;default:null" json:"to_ref_id"`
	ToRefName             *string        `gorm:"column:to_ref_name;default:null" json:"to_ref_name"`
	UserID                *int32         `gorm:"column:user_id;default:null" json:"user_id"`
	UserIDSubtitute       *int32         `gorm:"column:user_id_subtitute;default:null" json:"user_id_subtitute"`
}

// TableName Kunjungan's table name
func (*Kunjungan) TableName() string {
	return TableNameKunjungan
}
