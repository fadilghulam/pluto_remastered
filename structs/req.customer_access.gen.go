package structs

const TableNameCustomerAccess = "customer_access"

// CustomerAccess mapped from table <customer_access>
type CustomerAccess struct {
	ID          FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID  FlexibleString `gorm:"column:customer_id;not null" json:"customerIds"`
	AccessType  string         `gorm:"column:access_type;not null;default:DOUBLE CREDIT" json:"access_type"`
	DateStart   string         `gorm:"column:date_start;not null;default:CURRENT_DATE" json:"startDate"`
	DateEnd     string         `gorm:"column:date_end;not null;default:CURRENT_DATE" json:"endDate"`
	RequestAt   string         `gorm:"column:request_at;not null;default:now()" json:"request_at"`
	RequestedID int64          `gorm:"column:requested_id;not null" json:"employeeId"`
	IsApprove   *int16         `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt   *string        `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID   *int64         `gorm:"column:approve_id;default:null" json:"approve_id"`
	Note        *string        `gorm:"column:note;default:null" json:"note"`
	Attachment  *string        `gorm:"column:attachment;default:null" json:"attachment"`
	CreatedAt   string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt   string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	UserID      int32          `gorm:"column:user_id;not null" json:"userId"`
}

// TableName CustomerAccess's table name
func (*CustomerAccess) TableName() string {
	return TableNameCustomerAccess
}
