package structs

const TableNamePenjualanDetail = "penjualan_detail"

// PenjualanDetail mapped from table <penjualan_detail>
type PenjualanDetail struct {
	ID                  FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	PenjualanID         FlexibleString `gorm:"column:penjualan_id;not null" json:"penjualan_id"`
	ProdukID            int32          `gorm:"column:produk_id;not null" json:"produk_id"`
	IDInventorySalesman *int64         `gorm:"column:id_inventory_salesman;default:null" json:"id_inventory_salesman"`
	NoBatch             *string        `gorm:"column:no_batch;default:null" json:"no_batch"`
	Jumlah              int32          `gorm:"column:jumlah;not null" json:"jumlah"`
	DtmCrt              string         `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd              string         `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	Harga               float64        `gorm:"column:harga;not null" json:"harga"`
	Diskon              float64        `gorm:"column:diskon;not null" json:"diskon"`
	SyncKey             string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	Label               *string        `gorm:"column:label;default:null" json:"label"`
	Condition           *string        `gorm:"column:condition;default:null" json:"condition"`
	Pita                FlexibleString `gorm:"column:pita;default:null" json:"pita"`
	TeamleaderID        *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	UserID              *int32         `gorm:"column:user_id;default:null" json:"user_id"`
}

// TableName PenjualanDetail's table name
func (*PenjualanDetail) TableName() string {
	return TableNamePenjualanDetail
}
