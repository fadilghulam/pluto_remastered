package structs

const TableNameMdOutlet = "md.outlet"

// Md Outlet mapped from table <md.outlet>
type MdOutlet struct {
	ID                    FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name                  string         `gorm:"column:name;not null" json:"name"`
	OutletName            *string        `gorm:"column:outlet_name" json:"outlet_name"`
	Phone                 *string        `gorm:"column:phone" json:"phone"`
	Address               *string        `gorm:"column:address" json:"address"`
	Type                  int32          `gorm:"column:type;not null" json:"type"`
	LatitudeLongitude     *string        `gorm:"column:latitude_longitude" json:"latitude_longitude"`
	Provinsi              *string        `gorm:"column:provinsi" json:"provinsi"`
	Kabupaten             *string        `gorm:"column:kabupaten" json:"kabupaten"`
	Kecamatan             *string        `gorm:"column:kecamatan" json:"kecamatan"`
	Kelurahan             *string        `gorm:"column:kelurahan" json:"kelurahan"`
	OutletPhoto           *StringArray   `gorm:"column:outlet_photo" json:"outlet_photo"`
	SubjectTypeID         int16          `gorm:"column:subject_type_id;not null" json:"subject_type_id"`
	Note                  *string        `gorm:"column:note" json:"note"`
	MerchandiserIDCreator *int32         `gorm:"column:merchandiser_id_creator" json:"merchandiser_id_creator"`
	TeamleaderIDCreator   *int32         `gorm:"column:teamleader_id_creator" json:"teamleader_id_creator"`
	CreatedAt             string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt             string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	BranchID              int16          `gorm:"column:branch_id;not null" json:"branch_id"`
	RayonID               int16          `gorm:"column:rayon_id;not null" json:"rayon_id"`
	SrID                  int16          `gorm:"column:sr_id;not null" json:"sr_id"`
	IsConsume             *int16         `gorm:"column:is_consume" json:"is_consume"`
	EmployeeIDConsume     *int64         `gorm:"column:employee_id_consume" json:"employee_id_consume"`
	ConsumeAt             *string        `gorm:"column:consume_at" json:"consume_at"`
	SyncKey               string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
}

// TableName MdOutlet's table name
func (*MdOutlet) TableName() string {
	return TableNameMdOutlet
}
