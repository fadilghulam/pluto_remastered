package structs

const TableNameSalesmanAccessKunjungan = "salesman_access_kunjungan"

// SalesmanAccessKunjungan mapped from table <salesman_access_kunjungan>
type SalesmanAccessKunjungan struct {
	ID               FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	RequestAt        string         `gorm:"column:request_at;not null;default:now()" json:"request_at"`
	SalesmanID       *string        `gorm:"column:salesman_id;default:null" json:"salesmanIds"`
	DateStart        string         `gorm:"column:date_start;not null" json:"startDate"`
	DateEnd          *string        `gorm:"column:date_end;default:null" json:"endDate"`
	RequestedID      *int32         `gorm:"column:requested_id;default:null" json:"employeeId"`
	CreatedAt        string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt        string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	IsApprove        *int16         `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt        *string        `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID        *int32         `gorm:"column:approve_id;default:null" json:"approve_id"`
	Note             *string        `gorm:"column:note;default:null" json:"note"`
	Attachment       *string        `gorm:"column:attachment;default:null" json:"attachment"`
	MinMinute        *int32         `gorm:"column:min_minute;default:null" json:"visitMin"`
	MaxMinute        *int32         `gorm:"column:max_minute;default:null" json:"visitMax"`
	MaxRadiusOutlet  *int32         `gorm:"column:max_radius_outlet;default:null" json:"outletRadius"`
	MaxRadiusCheckIn *int32         `gorm:"column:max_radius_check_in;default:null" json:"checkOutRadius"`
	UserID           int32          `gorm:"column:user_id;not null" json:"userId"`
}

// TableName SalesmanAccessKunjungan's table name
func (*SalesmanAccessKunjungan) TableName() string {
	return TableNameSalesmanAccessKunjungan
}
