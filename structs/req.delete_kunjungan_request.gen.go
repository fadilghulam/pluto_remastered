package structs

const TableNameDeleteKunjunganRequest = "delete_kunjungan_request"

// DeleteKunjunganRequest mapped from table <delete_kunjungan_request>
type DeleteKunjunganRequest struct {
	EmployeeID  int64          `gorm:"column:employee_id;not null" json:"employeeId"`
	Datetime    string         `gorm:"column:datetime;not null;default:now()" json:"datetime"`
	KunjunganID FlexibleString `gorm:"column:kunjungan_id;not null" json:"kunjunganId"`
	Note        *string        `gorm:"column:note;default:null" json:"note"`
	IsApprove   *int16         `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt   *string        `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID   *int64         `gorm:"column:approve_id;default:null" json:"approve_id"`
	CreatedAt   string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt   string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID      int32          `gorm:"column:user_id;not null" json:"userId"`
}

// TableName DeleteKunjunganRequest's table name
func (*DeleteKunjunganRequest) TableName() string {
	return TableNameDeleteKunjunganRequest
}
