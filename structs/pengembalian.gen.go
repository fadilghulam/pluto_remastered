package structs

const TableNamePengembalian = "pengembalian"

// Pengembalian mapped from table <pengembalian>
type Pengembalian struct {
	ID                  FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	PenjualanID         FlexibleString `gorm:"column:penjualan_id;not null" json:"penjualan_id"`
	TanggalPengembalian string         `gorm:"column:tanggal_pengembalian;not null" json:"tanggal_pengembalian"`
	DtmCrt              string         `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd              string         `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	ImageNota           *string        `gorm:"column:image_nota;default:null" json:"image_nota"`
	IsKredit            int32          `gorm:"column:is_kredit;not null" json:"is_kredit"`
	SalesmanID          *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	SyncKey             string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	CustomerID          *string        `gorm:"column:customer_id;default:null" json:"customer_id"`
	BranchIDOld         *int16         `gorm:"column:branch_id_old;default:null" json:"branch_id_old"`
	BranchID            *int16         `gorm:"column:branch_id;default:null" json:"branch_id"`
	SrID                *int16         `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID             *int16         `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	AreaID              *int32         `gorm:"column:area_id;default:null" json:"area_id"`
	SalesmanTipe        *string        `gorm:"column:salesman_tipe;default:null" json:"salesman_tipe"`
	LatitudeLongitude   *string        `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
	Note                *string        `gorm:"column:note;default:null" json:"note"`
	CustomerTipe        *string        `gorm:"column:customer_tipe;default:null" json:"customer_tipe"`
	CustomerTypeID      *int16         `gorm:"column:customer_type_id;default:null" json:"customer_type_id"`
	SalesmanTypeID      *int16         `gorm:"column:salesman_type_id;default:null" json:"salesman_type_id"`
	MerchandiserID      *int32         `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	ToRefID             *int32         `gorm:"column:to_ref_id;default:null" json:"to_ref_id"`
	ToRefName           *string        `gorm:"column:to_ref_name;default:null" json:"to_ref_name"`
	TeamleaderID        *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	UserID              *int32         `gorm:"column:user_id;default:null" json:"user_id"`
	UserIDSubtitute     *int32         `gorm:"column:user_id_subtitute;default:null" json:"user_id_subtitute"`
}

// TableName Pengembalian's table name
func (*Pengembalian) TableName() string {
	return TableNamePengembalian
}
