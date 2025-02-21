package structs

const TableNameRuteMoveRequest = "rute_move_request"

// RuteMoveRequest mapped from table <rute_move_request>
type RuteMoveRequest struct {
	ID                    FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	RequestedID           int32          `gorm:"column:requested_id;not null" json:"employeeId"`
	RequestAt             string         `gorm:"column:request_at;not null;default:now()" json:"request_at"`
	RuteID                string         `gorm:"column:rute_id;not null" json:"ruteIds"`
	CustomerID            FlexibleString `gorm:"column:customer_id;default:null" json:"customer_id"`
	SalesmanIDDestination int32          `gorm:"column:salesman_id_destination;not null" json:"toSalesmanId"`
	DateEffective         string         `gorm:"column:date_effective;not null" json:"accessDate"`
	Note                  *string        `gorm:"column:note;default:null" json:"note"`
	Attachment            *string        `gorm:"column:attachment;default:null" json:"attachment"`
	IsApprove             *int16         `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt             *string        `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID             *int32         `gorm:"column:approve_id;default:null" json:"approve_id"`
	CreatedAt             string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt             string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	UserID                int32          `gorm:"column:user_id;not null" json:"userId"`
}

// TableName RuteMoveRequest's table name
func (*RuteMoveRequest) TableName() string {
	return TableNameRuteMoveRequest
}
