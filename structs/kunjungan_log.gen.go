package structs

const TableNameKunjunganLog = "kunjungan_log"

// KunjunganLog mapped from table <kunjungan_log>
type KunjunganLog struct {
	ID                  FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	SalesmanID          *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	CustomerID          FlexibleString `gorm:"column:customer_id;default:null" json:"customer_id"`
	CheckinAt           *string        `gorm:"column:checkin_at;default:null" json:"checkin_at"`
	CheckoutAt          *string        `gorm:"column:checkout_at;default:null" json:"checkout_at"`
	LatlongIn           *string        `gorm:"column:latlong_in;default:null" json:"latlong_in"`
	LatlongOut          *string        `gorm:"column:latlong_out;default:null" json:"latlong_out"`
	Photo               *string        `gorm:"column:photo;default:null" json:"photo"`
	PhotoIn             *string        `gorm:"column:photo_in;default:null" json:"photo_in"`
	PhotoOut            *string        `gorm:"column:photo_out;default:null" json:"photo_out"`
	DistanceIn          *int64         `gorm:"column:distance_in;default:null" json:"distance_in"`
	DistanceOut         *int64         `gorm:"column:distance_out;default:null" json:"distance_out"`
	AccuracyIn          *int16         `gorm:"column:accuracy_in;default:null" json:"accuracy_in"`
	AccuracyOut         *int16         `gorm:"column:accuracy_out;default:null" json:"accuracy_out"`
	IsMockLocationIn    *int16         `gorm:"column:is_mock_location_in;default:null" json:"is_mock_location_in"`
	IsMockLocationOut   *int16         `gorm:"column:is_mock_location_out;default:null" json:"is_mock_location_out"`
	VisitStatusID       *int16         `gorm:"column:visit_status_id;default:null" json:"visit_status_id"`
	VisitStatusNote     *string        `gorm:"column:visit_status_note;default:null" json:"visit_status_note"`
	IsVisitStatusAction *int16         `gorm:"column:is_visit_status_action;default:null" json:"is_visit_status_action"`
	VisitStatusActionAt *string        `gorm:"column:visit_status_action_at;default:null" json:"visit_status_action_at"`
	RefuseBuy           *string        `gorm:"column:refuse_buy;default:null" json:"refuse_buy"`
	RefusePay           *string        `gorm:"column:refuse_pay;default:null" json:"refuse_pay"`
	SrID                *int16         `gorm:"column:sr_id;default:null" json:"sr_id"`
	RayonID             *int16         `gorm:"column:rayon_id;default:null" json:"rayon_id"`
	BranchID            *int16         `gorm:"column:branch_id;default:null" json:"branch_id"`
	AreaID              *int32         `gorm:"column:area_id;default:null" json:"area_id"`
	CreatedAt           string         `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt           string         `gorm:"column:updated_at;default:now()" json:"updated_at"`
	CustomerTipe        *string        `gorm:"column:customer_tipe;default:null" json:"customer_tipe"`
	SyncKey             string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	CustomerTypeID      *int16         `gorm:"column:customer_type_id;default:null" json:"customer_type_id"`
	SalesmanTypeID      *int16         `gorm:"column:salesman_type_id;default:null" json:"salesman_type_id"`
	MerchandiserID      *int32         `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	TeamleaderID        *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	SpvID               *int32         `gorm:"column:spv_id;default:null" json:"spv_id"`
	Provinsi            *string        `gorm:"column:provinsi;default:null" json:"provinsi"`
	Kabupaten           *string        `gorm:"column:kabupaten;default:null" json:"kabupaten"`
	Kecamatan           *string        `gorm:"column:kecamatan;default:null" json:"kecamatan"`
	Kelurahan           *string        `gorm:"column:kelurahan;default:null" json:"kelurahan"`
	CustomerName        *string        `gorm:"column:customer_name;default:null" json:"customer_name"`
	OutletID            FlexibleString `gorm:"column:outlet_id;default:null" json:"outlet_id"`
	SubjectTypeID       *int32         `gorm:"column:subject_type_id;default:null" json:"subject_type_id"`
	OutletTypeID        *int32         `gorm:"column:outlet_type_id;default:null" json:"outlet_type_id"`
	ToRefID             *int32         `gorm:"column:to_ref_id;default:null" json:"to_ref_id"`
	ToRefName           *string        `gorm:"column:to_ref_name;default:null" json:"to_ref_name"`
	UserID              *int32         `gorm:"column:user_id;default:null" json:"user_id"`
	UserIDSubtitute     *int32         `gorm:"column:user_id_subtitute;default:null" json:"user_id_subtitute"`
}

// TableName KunjunganLog's table name
func (*KunjunganLog) TableName() string {
	return TableNameKunjunganLog
}
