package structs

const TableNameStokGudangPengiriman = "stok_gudang_pengiriman"

// StokGudangPengiriman mapped from table <stok_gudang_pengiriman>
type StokGudangPengiriman struct {
	ID                    int64   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AdminGudangPengirimID int16   `gorm:"column:admin_gudang_pengirim_id;not null" json:"admin_gudang_pengirim_id"`
	AdminGudangPenerimaID *int16  `gorm:"column:admin_gudang_penerima_id" json:"admin_gudang_penerima_id"`
	ConditionSource       string  `gorm:"column:condition_source;not null" json:"condition_source"`
	ConditionDestination  string  `gorm:"column:condition_destination;not null" json:"condition_destination"`
	TanggalKirim          string  `gorm:"column:tanggal_kirim;not null" json:"tanggal_kirim"`
	TanggalTerima         *string `gorm:"column:tanggal_terima" json:"tanggal_terima"`
	Tag                   *string `gorm:"column:tag;comment:SHIPMENT, ADJUSTMENT, SWITCH CONDITION, REJECT" json:"tag"` // SHIPMENT, ADJUSTMENT, SWITCH CONDITION, REJECT
	SuratJalan            *string `gorm:"column:surat_jalan" json:"surat_jalan"`
	Lampiran              *string `gorm:"column:lampiran" json:"lampiran"`
	DtmCrt                string  `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd                string  `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	Status                *string `gorm:"column:status;comment:IN TRANSIT, ACCEPTED,  WAITING CONFIRMATION, REJECTED" json:"status"` // IN TRANSIT, ACCEPTED,  WAITING CONFIRMATION, REJECTED
	GudangNamaAsalOld     *string `gorm:"column:gudang_nama_asal_old" json:"gudang_nama_asal_old"`
	GudangNamaTujuanOld   *string `gorm:"column:gudang_nama_tujuan_old" json:"gudang_nama_tujuan_old"`
	SuratJalanVendor      *string `gorm:"column:surat_jalan_vendor" json:"surat_jalan_vendor"`
	CatatanPenerima       *string `gorm:"column:catatan_penerima" json:"catatan_penerima"`
	CatatanPengirim       *string `gorm:"column:catatan_pengirim" json:"catatan_pengirim"`
	GudangIDAsal          *int16  `gorm:"column:gudang_id_asal" json:"gudang_id_asal"`
	GudangIDTujuan        *int16  `gorm:"column:gudang_id_tujuan" json:"gudang_id_tujuan"`
	UpdatePriceAt         *string `gorm:"column:update_price_at" json:"update_price_at"`
}

// TableName StokGudangPengiriman's table name
func (*StokGudangPengiriman) TableName() string {
	return TableNameStokGudangPengiriman
}
