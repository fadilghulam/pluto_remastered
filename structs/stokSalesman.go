package structs

const TableNameStokSalesman = "public.stok_salesman"

// Stok Salesman mapped from table <stok_salesman>
type StokSalesman struct {
	ID              int64   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	StokGudangId    *int32  `gorm:"column:stok_gudang_id" json:"stok_gudang_id"`
	SalesmanId      *int32  `gorm:"column:salesman_id" json:"salesman_id"`
	ProdukId        *int16  `gorm:"column:produk_id" json:"produk_id"`
	StokAwal        *int32  `gorm:"column:stok_awal" json:"stok_awal"`
	StokAkhir       *int32  `gorm:"column:stok_akhir" json:"stok_akhir"`
	TanggalStok     *string `gorm:"column:tanggal_stok" json:"tanggal_stok"`
	DtmCrt          string  `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd          string  `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	ConfirmKey      *string `gorm:"column:confirm_key" json:"confirm_key"`
	IsComplete      int16   `gorm:"column:is_complete;default: 0" json:"is_complete"`
	TanggalSo       *string `gorm:"column:tanggal_so" json:"tanggal_so"`
	SoAdminGudangId *int16  `gorm:"column:so_admin_gudang_id" json:"so_admin_gudang_id"`
	Condition       *string `gorm:"column:condition" json:"condition"`
	Pita            *int32  `gorm:"column:pita" json:"pita"`
	MerchandiserId  *int32  `gorm:"column:merchandiser_id" json:"merchandiser_id"`
	UserId          *int32  `gorm:"column:user_id" json:"user_id"`
	UserIdSubtitute *int32  `gorm:"column:user_id_subtitute" json:"user_id_subtitute"`
	StokUserId      *int64  `gorm:"column:stok_user_id" json:"stok_user_id"`
}

// TableName Stok Salesman's table name
func (*StokSalesman) TableName() string {
	return TableNameStokSalesman
}
