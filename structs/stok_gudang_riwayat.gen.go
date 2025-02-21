package structs

import "time"

const TableNameStokGudangRiwayat = "stok_gudang_riwayat"

// StokGudangRiwayat mapped from table <stok_gudang_riwayat>
type StokGudangRiwayat struct {
	ID                     int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AdminGudangID          *int32    `gorm:"column:admin_gudang_id" json:"admin_gudang_id"`
	StokGudangID           *int32    `gorm:"column:stok_gudang_id" json:"stok_gudang_id"`
	Jumlah                 int32     `gorm:"column:jumlah;not null" json:"jumlah"`
	TanggalRiwayat         string    `gorm:"column:tanggal_riwayat;default:now()" json:"tanggal_riwayat"`
	Aksi                   *string   `gorm:"column:aksi" json:"aksi"`
	NoSurat                *string   `gorm:"column:no_surat" json:"no_surat"`
	FotoSurat              *string   `gorm:"column:foto_surat" json:"foto_surat"`
	DtmCrt                 time.Time `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd                 time.Time `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	HargaBeli              *int32    `gorm:"column:harga_beli" json:"harga_beli"`
	ProdukID               int16     `gorm:"column:produk_id;not null" json:"produk_id"`
	IsValidate             int32     `gorm:"column:is_validate;not null" json:"is_validate"`
	Catatan                *string   `gorm:"column:catatan" json:"catatan"`
	CatatanUmum            *string   `gorm:"column:catatan_umum" json:"catatan_umum"`
	Condition              *string   `gorm:"column:condition" json:"condition"`
	Pita                   *string   `gorm:"column:pita" json:"pita"`
	StokGudangPengirimanID *int64    `gorm:"column:stok_gudang_pengiriman_id" json:"stok_gudang_pengiriman_id"`
}

// TableName StokGudangRiwayat's table name
func (*StokGudangRiwayat) TableName() string {
	return TableNameStokGudangRiwayat
}
