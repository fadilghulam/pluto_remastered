package structs

const TableNameCustomerTypeRequest = "customer_type_request"

// CustomerTypeRequest mapped from table <customer_type_request>
type CustomerTypeRequest struct {
	ID             FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerTypeID int16          `gorm:"column:customer_type_id;not null" json:"customerTypeId"`
	CustomerID     FlexibleString `gorm:"column:customer_id;not null" json:"customerIds"`
	DateEffective  string         `gorm:"column:date_effective;not null;default:now()" json:"accessDate"`
	RequestAt      string         `gorm:"column:request_at;not null;default:now()" json:"request_at"`
	RequestedID    *int64         `gorm:"column:requested_id;default:null" json:"employeeId"`
	IsApprove      *int16         `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt      *string        `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID      *int64         `gorm:"column:approve_id;default:null" json:"approve_id"`
	Note           *string        `gorm:"column:note;default:null" json:"note"`
	Attachment     *string        `gorm:"column:attachment;default:null" json:"attachment"`
	CreatedAt      string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt      string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	UserID         int32          `gorm:"column:user_id;not null" json:"userId"`
}

// TableName CustomerTypeRequest's table name
func (*CustomerTypeRequest) TableName() string {
	return TableNameCustomerTypeRequest
}
