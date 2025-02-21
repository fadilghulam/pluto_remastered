package structs

const TableNamePengembalianDetail = "pengembalian_detail"

// PengembalianDetail mapped from table <pengembalian_detail>
type PengembalianDetail struct {
	ID                  FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	PengembalianID      FlexibleString `gorm:"column:pengembalian_id;not null" json:"pengembalian_id"`
	ProdukID            int32          `gorm:"column:produk_id;not null" json:"produk_id"`
	IDInventorySalesman *int64         `gorm:"column:id_inventory_salesman;default:null" json:"id_inventory_salesman"`
	NoBatch             *string        `gorm:"column:no_batch;default:null" json:"no_batch"`
	Jumlah              int32          `gorm:"column:jumlah;not null" json:"jumlah"`
	DtmCrt              string         `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd              string         `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	Harga               float64        `gorm:"column:harga;not null;default:0" json:"harga"`
	SyncKey             string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	Condition           *string        `gorm:"column:condition;default:null" json:"condition"`
	Pita                *string        `gorm:"column:pita;default:null" json:"pita"`
	TeamleaderID        *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	UserID              *int32         `gorm:"column:user_id;default:null" json:"user_id"`
}

// TableName PengembalianDetail's table name
func (*PengembalianDetail) TableName() string {
	return TableNamePengembalianDetail
}
