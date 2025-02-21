package structs

const TableNameSalesmanRequestSo = "salesman_request_so"

// SalesmanRequestSo mapped from table <salesman_request_so>
type SalesmanRequestSo struct {
	ID          FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	SalesmanID  *int32         `gorm:"column:salesman_id;default:null" json:"salesmanIds"`
	Date        *string        `gorm:"column:date;default:null" json:"accessDate"`
	Type        string         `gorm:"column:type;default:so" json:"type"`
	CreatedAt   string         `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt   string         `gorm:"column:updated_at;default:now()" json:"updated_at"`
	RequestedID *int64         `gorm:"column:requested_id;default:null" json:"employeeId"`
	IsApprove   *int16         `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt   *string        `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID   *int32         `gorm:"column:approve_id;default:null" json:"approve_id"`
	RequestAt   *string        `gorm:"column:request_at;default:null" json:"request_at"`
	UserID      int32          `gorm:"column:user_id;not null" json:"userId"`
}

// TableName SalesmanRequestSo's table name
func (*SalesmanRequestSo) TableName() string {
	return TableNameSalesmanRequestSo
}
