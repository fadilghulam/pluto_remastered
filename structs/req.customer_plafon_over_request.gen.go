package structs

const TableNameCustomerPlafonOverRequest = "customer_plafon_over_request"

// CustomerPlafonOverRequest mapped from table <customer_plafon_over_request>
type CustomerPlafonOverRequest struct {
	ID          FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	RequestAt   string         `gorm:"column:request_at;not null;default:now()" json:"request_at"`
	RequestedID *int64         `gorm:"column:requested_id;default:null" json:"requested_id"`
	IsApprove   *int16         `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt   *string        `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID   *int64         `gorm:"column:approve_id;default:null" json:"approve_id"`
	Note        *string        `gorm:"column:note;default:null" json:"note"`
	Attachment  *string        `gorm:"column:attachment;default:null" json:"attachment"`
	CreatedAt   string         `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt   string         `gorm:"column:updated_at;default:now()" json:"updated_at"`
	DateStart   string         `gorm:"column:date_start;not null;default:now()" json:"date_start"`
	DateEnd     string         `gorm:"column:date_end;not null;default:now()" json:"date_end"`
	ExecutedAt  *string        `gorm:"column:executed_at;default:null" json:"executed_at"`
	UserID      int32          `gorm:"column:user_id;not null" json:"userId"`
}

// TableName CustomerPlafonOverRequest's table name
func (*CustomerPlafonOverRequest) TableName() string {
	return TableNameCustomerPlafonOverRequest
}
