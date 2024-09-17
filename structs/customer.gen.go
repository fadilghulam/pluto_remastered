package structs

import (
	"time"
)

const TableNameCustomer = "customer"

// Customer mapped from table <customer>
type Customer struct {
	ID                      int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	SalesmanID              int32     `gorm:"column:salesman_id;not null" json:"salesman_id"`
	Name                    string    `gorm:"column:name;not null" json:"name"`
	OutletName              string    `gorm:"column:outlet_name" json:"outlet_name"`
	Alamat                  string    `gorm:"column:alamat" json:"alamat"`
	Phone                   string    `gorm:"column:phone" json:"phone"`
	Tipe                    int32     `gorm:"column:tipe;not null" json:"tipe"`
	LatitudeLongitude       string    `gorm:"column:latitude_longitude" json:"latitude_longitude"`
	ImageKtp                string    `gorm:"column:image_ktp" json:"image_ktp"`
	ImageToko               string    `gorm:"column:image_toko" json:"image_toko"`
	IsAcc                   int32     `gorm:"column:is_acc;not null" json:"is_acc"`
	IsAktif                 int32     `gorm:"column:is_aktif;not null;default:1" json:"is_aktif"`
	DtmCrt                  time.Time `gorm:"column:dtm_crt;not null;default:now()" json:"dtm_crt"`
	DtmUpd                  time.Time `gorm:"column:dtm_upd;not null;default:now()" json:"dtm_upd"`
	Plafond                 float64   `gorm:"column:plafond;not null" json:"plafond"`
	Piutang                 float64   `gorm:"column:piutang;not null" json:"piutang"`
	KodeCustomer            string    `gorm:"column:kode_customer;not null" json:"kode_customer"`
	QrCode                  string    `gorm:"column:qr_code" json:"qr_code"`
	Nik                     string    `gorm:"column:nik;default:-" json:"nik"`
	Diskon                  float64   `gorm:"column:diskon;not null" json:"diskon"`
	Provinsi                string    `gorm:"column:provinsi" json:"provinsi"`
	Kabupaten               string    `gorm:"column:kabupaten" json:"kabupaten"`
	Kecamatan               string    `gorm:"column:kecamatan" json:"kecamatan"`
	Kelurahan               string    `gorm:"column:kelurahan" json:"kelurahan"`
	KawasanToko             string    `gorm:"column:kawasan_toko" json:"kawasan_toko"`
	HariKunjungan           string    `gorm:"column:hari_kunjungan" json:"hari_kunjungan"`
	FrekKunjungan           string    `gorm:"column:frek_kunjungan" json:"frek_kunjungan"`
	KawasanTokoOth          string    `gorm:"column:kawasan_toko_oth" json:"kawasan_toko_oth"`
	FrekKunjunganOth        string    `gorm:"column:frek_kunjungan_oth" json:"frek_kunjungan_oth"`
	IsVerifikasi            int32     `gorm:"column:is_verifikasi;not null;comment:digunakan untuk membedakan new outlet disetujui" json:"is_verifikasi"` // digunakan untuk membedakan new outlet disetujui
	ImageTokoAfter          string    `gorm:"column:image_toko_after" json:"image_toko_after"`
	Validated               int32     `gorm:"column:validated" json:"validated"`
	SyncKey                 string    `gorm:"column:sync_key" json:"sync_key"`
	AksesDoubleKredit       time.Time `gorm:"column:akses_double_kredit" json:"akses_double_kredit"`
	TanggalVerifikasi       time.Time `gorm:"column:tanggal_verifikasi" json:"tanggal_verifikasi"`
	VerifiedBy              int32     `gorm:"column:verified_by" json:"verified_by"`
	IsVerifikasiLokasi      int16     `gorm:"column:is_verifikasi_lokasi" json:"is_verifikasi_lokasi"`
	VisitExtra              time.Time `gorm:"column:visit_extra" json:"visit_extra"`
	IsMandiri               int16     `gorm:"column:is_mandiri;not null" json:"is_mandiri"`
	SetMandiriBy            int32     `gorm:"column:set_mandiri_by;comment:diset oleh user id pada tabel user" json:"set_mandiri_by"` // diset oleh user id pada tabel user
	IsKasus                 int16     `gorm:"column:is_kasus;not null;comment:customer kasus..." json:"is_kasus"`                     // customer kasus...
	SalesmanTemp            int64     `gorm:"column:salesman_temp" json:"salesman_temp"`
	SisaKreditNoo           int16     `gorm:"column:sisa_kredit_noo;not null" json:"sisa_kredit_noo"`
	AreaID                  int32     `gorm:"column:area_id" json:"area_id"`
	SalesmanTypeID          int16     `gorm:"column:salesman_type_id" json:"salesman_type_id"`
	IsHandover              int16     `gorm:"column:is_handover" json:"is_handover"`
	CreatedID               int32     `gorm:"column:created_id;comment:employee_id" json:"created_id"` // employee_id
	DateLastTransaction     time.Time `gorm:"column:date_last_transaction" json:"date_last_transaction"`
	DateLastVisitBySalesman time.Time `gorm:"column:date_last_visit_by_salesman" json:"date_last_visit_by_salesman"`
	OutletPhoto             string    `gorm:"column:outlet_photo" json:"outlet_photo"`
	Note                    string    `gorm:"column:note" json:"note"`
	SubjectTypeID           int16     `gorm:"column:subject_type_id" json:"subject_type_id"`
	SalesmanIDCreator       int32     `gorm:"column:salesman_id_creator" json:"salesman_id_creator"`
	MerchandiserIDCreator   int32     `gorm:"column:merchandiser_id_creator" json:"merchandiser_id_creator"`
	TeamleaderIDCreator     int32     `gorm:"column:teamleader_id_creator" json:"teamleader_id_creator"`
	BranchID                int16     `gorm:"column:branch_id" json:"branch_id"`
	RayonID                 int16     `gorm:"column:rayon_id" json:"rayon_id"`
	SrID                    int16     `gorm:"column:sr_id" json:"sr_id"`
	IsConsume               int16     `gorm:"column:is_consume" json:"is_consume"`
	EmployeeIDConsume       int64     `gorm:"column:employee_id_consume" json:"employee_id_consume"`
	ConsumeAt               time.Time `gorm:"column:consume_at" json:"consume_at"`
	/*
		REGULAR
		RECOMENDATION
	*/
	Tag string `gorm:"column:tag;not null;default:REGULAR;comment:REGULAR\nRECOMENDATION" json:"tag"`
}

// TableName Customer's table name
func (*Customer) TableName() string {
	return TableNameCustomer
}
