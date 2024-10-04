package structs

import (
	"time"
)

const TableNameStokSalesmanRiwayat = "public.stok_salesman_riwayat"

// Salesman mapped from table <salesman>
type StokSalesmanRiwayat struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	StokSalesmanId  int64     `gorm:"column:stok_salesman_id;default:null" json:"stok_salesman_id"`
	AdminGudangId   int16     `gorm:"column:admin_gudang_id;default:null" json:"admin_gudang_id"`
	ProdukId        int16     `gorm:"column:produk_id;not null" json:"produk_id"`
	Jumlah          int32     `gorm:"column:jumlah;not null" json:"jumlah"`
	Aksi            string    `gorm:"column:aksi;not null" json:"aksi"`
	TanggalRiwayat  time.Time `gorm:"column:tanggal_riwayat;not null" json:"tanggal_riwayat"`
	DtmCrt          time.Time `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd          time.Time `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	ConfirmKey      string    `gorm:"column:confirm_key;default:null" json:"confirm_key"`
	IsValidate      int16     `gorm:"column:is_validate;default: 0" json:"is_validate"`
	SalesmanId      int32     `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	GudangNamaOld   string    `gorm:"column:gudang_nama_old;default:null" json:"gudang_nama_old"`
	SerahTerima     string    `gorm:"column:serah_terima;default:null" json:"serah_terima"`
	Condition       string    `gorm:"column:condition" json:"condition"`
	Pita            string    `gorm:"column:pita" json:"pita"`
	GudangId        int16     `gorm:"column:gudang_id" json:"gudang_id"`
	ParentId        int64     `gorm:"column:parent_id" json:"parent_id"`
	MerchandiserId  int32     `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	TeamleaderId    int32     `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	UserId          int32     `gorm:"column:user_id;default:null" json:"user_id"`
	UserIdSubtitute int32     `gorm:"column:user_id_subtitute;default:null" json:"user_id_subtitute"`
}

// TableName Salesman's table name
func (*StokSalesmanRiwayat) TableName() string {
	return TableNameStokSalesmanRiwayat
}
