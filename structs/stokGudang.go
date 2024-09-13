package structs

import (
	"time"
)

const TableNameStokGudang = "public.stok_gudang"

type StokGudang struct {
	ID            int16     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	BranchIdOld   int16     `gorm:"column:branch_id_old;default:null" json:"branch_id_old"`
	ProdukId      int16     `gorm:"column:produk_id" json:"produk_id"`
	Harga         int32     `gorm:"column:harga" json:"harga"`
	Jumlah        int64     `gorm:"column:jumlah" json:"jumlah"`
	DtmCrt        time.Time `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd        time.Time `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	Batch         string    `gorm:"column:batch;default:null" json:"batch"`
	GudangNamaOld string    `gorm:"column:gudang_nama_old;default:null" json:"gudang_nama_old"`
	Condition     string    `gorm:"column:condition" json:"condition"`
	Pita          string    `gorm:"column:pita" json:"pita"`
	GudangId      *int16    `gorm:"column:gudang_id" json:"gudang_id"`
}

// TableName Salesman's table name
func (*StokGudang) TableName() string {
	return TableNameStokGudang
}
