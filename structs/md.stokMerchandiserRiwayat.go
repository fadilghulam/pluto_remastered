package structs

import (
	"time"
)

const TableNameStokMerchandiserRiwayat = "md.stok_merchandiser_riwayat"

type StokMerchandiserRiwayat struct {
	ID                 int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	StokMerchandiserId int64     `gorm:"column:stok_merchandiser_id;default:null" json:"stok_merchandiser_id"`
	AdminGudangId      int16     `gorm:"column:admin_gudang_id;default:null" json:"admin_gudang_id"`
	ItemId             int16     `gorm:"column:item_id;not null" json:"item_id"`
	Jumlah             int32     `gorm:"column:jumlah;not null" json:"jumlah"`
	Aksi               string    `gorm:"column:aksi;not null" json:"aksi"`
	TanggalRiwayat     time.Time `gorm:"column:tanggal_riwayat;not null" json:"tanggal_riwayat"`
	DtmCrt             time.Time `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd             time.Time `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	ConfirmKey         string    `gorm:"column:confirm_key;default:null" json:"confirm_key"`
	IsValidate         int16     `gorm:"column:is_validate;default: 0" json:"is_validate"`
	MerchandiserId     int32     `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	GudangNamaOld      string    `gorm:"column:gudang_nama_old;default:null" json:"gudang_nama_old"`
	SerahTerima        string    `gorm:"column:serah_terima;default:null" json:"serah_terima"`
	GudangId           int16     `gorm:"column:gudang_id" json:"gudang_id"`
	ParentId           int64     `gorm:"column:parent_id" json:"parent_id"`
	UserId             int32     `gorm:"column:user_id;default:null" json:"user_id"`
	UserIdSubtitute    int32     `gorm:"column:user_id_subtitute;default:null" json:"user_id_subtitute"`
}

// TableName Salesman's table name
func (*StokMerchandiserRiwayat) TableName() string {
	return TableNameStokMerchandiserRiwayat
}
