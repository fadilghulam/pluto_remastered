package structs

const TableNameSalesman = "public.salesman"

// Salesman mapped from table <salesman>
type Salesman struct {
	ID                   int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	BranchIDOld          *int32     `gorm:"column:branch_id_old" json:"branch_id_old"`
	Name                 string     `gorm:"column:name;not null" json:"name"`
	Phone                string     `gorm:"column:phone;not null" json:"phone"`
	Email                string     `gorm:"column:email;not null" json:"email"`
	IsAktif              int32      `gorm:"column:is_aktif;not null" json:"is_aktif"`
	DtmCrt               string     `gorm:"column:dtm_crt;not null;default:(now() + '12:00:00" json:"dtm_crt"`
	DtmUpd               string     `gorm:"column:dtm_upd;not null;default:(now() + '12:00:00" json:"dtm_upd"`
	TeamleaderID         *int32     `gorm:"column:teamleader_id" json:"teamleader_id"`
	TipeSalesman         string     `gorm:"column:tipe_salesman;default:MOTORIS" json:"tipe_salesman"`
	AksesRetur           *string    `gorm:"column:akses_retur" json:"akses_retur"`
	AksesKredit          *string    `gorm:"column:akses_kredit" json:"akses_kredit"`
	IsFixedRoute         int16      `gorm:"column:is_fixed_route;not null" json:"is_fixed_route"`
	IsKreditRetail       int16      `gorm:"column:is_kredit_retail;not null" json:"is_kredit_retail"`
	IsKreditSubGrosir    int16      `gorm:"column:is_kredit_sub_grosir;not null" json:"is_kredit_sub_grosir"`
	IsKreditGrosir       int16      `gorm:"column:is_kredit_grosir;not null" json:"is_kredit_grosir"`
	AksesDoubleTransaksi *string    `gorm:"column:akses_double_transaksi" json:"akses_double_transaksi"`
	IsReturAll           int16      `gorm:"column:is_retur_all;not null" json:"is_retur_all"`
	BranchID             *int16     `gorm:"column:branch_id" json:"branch_id"`
	Nik                  *string    `gorm:"column:nik" json:"nik"`
	AksesKreditNoo       *string    `gorm:"column:akses_kredit_noo" json:"akses_kredit_noo"`
	UserType             string     `gorm:"column:user_type;not null;default:SALESMAN" json:"user_type"`
	AreaID               Int32Array `gorm:"column:area_id;type:int[]" json:"area_id"`
	UserID               *int32     `gorm:"column:user_id" json:"user_id"`
	SrID                 *int16     `gorm:"column:sr_id" json:"sr_id"`
	RayonID              *int16     `gorm:"column:rayon_id" json:"rayon_id"`
	SpvID                *int16     `gorm:"column:spv_id" json:"spv_id"`
	SalesmanTypeID       int16      `gorm:"column:salesman_type_id;not null;default:(1)" json:"salesman_type_id"`
	IsHaveCustomer       int16      `gorm:"column:is_have_customer;not null;default:1" json:"is_have_customer"`
	SkID                 *string    `gorm:"column:sk_id" json:"sk_id"`
	IgnorePlafonAccess   string     `gorm:"column:ignore_plafon_access;default:2024-12-31 20:00:00" json:"ignore_plafon_access"`
}

// TableName Salesman's table name
func (*Salesman) TableName() string {
	return TableNameSalesman
}
