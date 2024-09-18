package structs

import (
	"time"
)

const TableNameMdTransaction = "md.transaction"

// Md Transaction mapped from table <md.transaction>
type MdTransaction struct {
	ID                string      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	MerchandiserID    int32       `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	SalesmanID        int32       `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	SubjectTypeID     int16       `gorm:"column:subject_type_id;not null" json:"subject_type_id"`
	OutletID          string      `gorm:"column:outlet_id;default:null" json:"outlet_id"`
	Datetime          time.Time   `gorm:"column:datetime;not null" json:"datetime"`
	Photos            StringArray `gorm:"type:varchar[];column:photos;default:null" json:"photos"`
	PhotosBefore      StringArray `gorm:"type:varchar[];column:photos_before;default:null" json:"photos_before"`
	PhotosAfter       StringArray `gorm:"type:varchar[];column:photos_after;default:null" json:"photos_after"`
	TransactionTypeID int16       `gorm:"column:transaction_type_id;not null" json:"transaction_type_id"`
	Provinsi          string      `gorm:"column:provinsi;default:null" json:"provinsi"`
	ProvinsiID        int32       `gorm:"column:provinsi_id;default:null" json:"provinsi_id"`
	Kabupaten         string      `gorm:"column:kabupaten;default:null" json:"kabupaten"`
	KabupatenID       int32       `gorm:"column:kabupaten_id;default:null" json:"kabupaten_id"`
	Kecamatan         string      `gorm:"column:kecamatan;default:null" json:"kecamatan"`
	KecamatanID       int32       `gorm:"column:kecamatan_id;default:null" json:"kecamatan_id"`
	Kelurahan         string      `gorm:"column:kelurahan;default:null" json:"kelurahan"`
	KelurahanID       int64       `gorm:"column:kelurahan_id;default:null" json:"kelurahan_id"`
	SrID              int16       `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID           int16       `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	BranchID          int16       `gorm:"column:branch_id;default:null" json:"branch_id"`
	ProgramID         int32       `gorm:"column:program_id;default:null" json:"program_id"`
	EventID           int32       `gorm:"column:event_id;default:null" json:"event_id"`
	LatitudeLongitude string      `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
	CreatedAt         time.Time   `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time   `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	SyncKey           string      `gorm:"column:sync_key;default:now()" json:"sync_key"`
	ToRefID           int32       `gorm:"column:to_ref_id;default:null" json:"to_ref_id"`
	ToRefName         string      `gorm:"column:to_ref_name;default:null" json:"to_ref_name"`
	CustomerID        string      `gorm:"column:customer_id;default:null" json:"customer_id"`
	TransactionIDs    StringArray `gorm:"column:transaction_ids;default:null" json:"transaction_ids"`
	UserID            int32       `gorm:"column:user_id;default:null" json:"user_id"`
	UserIDSubtitute   int32       `gorm:"column:user_id_subtitute;default:null" json:"user_id_subtitute"`
}

// TableName MdTransaction's table name
func (*MdTransaction) TableName() string {
	return TableNameMdTransaction
}
