package structs

const TableNameSalesmanAccess = "salesman_access"

// SalesmanAccess mapped from table <salesman_access>
type SalesmanAccess struct {
	ID                FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	RequestAt         string         `gorm:"column:request_at;not null;default:now()" json:"request_at"`
	SalesmanID        string         `gorm:"column:salesman_id;not null" json:"salesmanIds"`
	AccessType        string         `gorm:"column:access_type;not null;default:CREDIT" json:"typeAccess"`
	DateStart         string         `gorm:"column:date_start;not null;default:now()" json:"startDate"`
	DateEnd           string         `gorm:"column:date_end;default:now()" json:"endDate"`
	RequestedID       int32          `gorm:"column:requested_id;not null" json:"employeeId"`
	CreatedAt         string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt         string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	IsApprove         *int16         `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt         *string        `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID         *int32         `gorm:"column:approve_id;default:null" json:"approve_id"`
	Note              *string        `gorm:"column:note;default:null" json:"note"`
	Attachment        *string        `gorm:"column:attachment;default:null" json:"attachment"`
	ActiveType        string         `gorm:"column:active_type;not null;default:BUKA" json:"typeOpen"`
	IsKreditRetail    *int16         `gorm:"column:is_kredit_retail;default:null" json:"isKreditRetail"`
	IsKreditSubGrosir *int16         `gorm:"column:is_kredit_sub_grosir;default:null" json:"isKreditSubGrosir"`
	IsKreditGrosir    *int16         `gorm:"column:is_kredit_grosir;default:null" json:"isKreditGrosir"`
	UserID            int32          `gorm:"column:user_id;not null" json:"userId"`
}

// TableName SalesmanAccess's table name
func (*SalesmanAccess) TableName() string {
	return TableNameSalesmanAccess
}
