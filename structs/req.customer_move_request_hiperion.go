package structs

const TableNameCustomerMoveRequestHiperion = "customer_move_request"

// CustomerMoveRequest mapped from table <customer_move_request>
type CustomerMoveRequestHiperion struct {
	ID                    FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	RequestedID           int64          `gorm:"column:requested_id;not null" json:"requestedId"`
	RequestAt             string         `gorm:"column:request_at;not null" json:"request_at"`
	SalesmanIDSource      *string        `gorm:"column:salesman_id_source;default:null" json:"salesSource"`
	SalesmanIDDestination *int32         `gorm:"column:salesman_id_destination;default:null" json:"salesTujuan"`
	CustomerID            FlexibleString `gorm:"column:customer_id;not null" json:"customerId"`
	DateEffective         string         `gorm:"column:date_effective;not null" json:"effectiveDate"`
	Attachment            *string        `gorm:"column:attachment;default:null" json:"attachment"`
	Note                  *string        `gorm:"column:note;default:null" json:"note"`
	Tag                   string         `gorm:"column:tag;not null" json:"tag"`
	IsApprove             *int16         `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt             *string        `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID             *int64         `gorm:"column:approve_id;default:null" json:"approve_id"`
	CreatedAt             string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt             string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

// TableName CustomerMoveRequest's table name
func (*CustomerMoveRequestHiperion) TableName() string {
	return TableNameCustomerMoveRequestHiperion
}
