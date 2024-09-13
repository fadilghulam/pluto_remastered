package structs

import (
	"time"
)

const TableNameStokGudangItem = "md.stok_gudang_item"

type StokGudangItem struct {
	ID          int16     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	BranchIdOld int16     `gorm:"column:branch_id_old;default:null" json:"branch_id_old"`
	ItemId      int16     `gorm:"column:item_id" json:"item_id"`
	Harga       int32     `gorm:"column:harga" json:"harga"`
	Jumlah      int64     `gorm:"column:jumlah" json:"jumlah"`
	DtmCrt      time.Time `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd      time.Time `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	Batch       string    `gorm:"column:batch;default:null" json:"batch"`
	GudangId    int16     `gorm:"column:gudang_id" json:"gudang_id"`
}

// TableName Salesman's table name
func (*StokGudangItem) TableName() string {
	return TableNameStokGudangItem
}
