package structs

const TableNameCustomerRelocation = "customer_relocation"

// CustomerRelocation mapped from table <customer_relocation>
type CustomerRelocation struct {
	ID                      FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	EmployeeID              int64          `gorm:"column:employee_id;not null" json:"employeeId"`
	CustomerID              FlexibleString `gorm:"column:customer_id;not null" json:"customerId"`
	RequestAt               string         `gorm:"column:request_at;not null;default:now()" json:"request_at"`
	LatitudeLongitudeBefore *string        `gorm:"column:latitude_longitude_before;default:null" json:"latlongBefore"`
	LatitudeLongitudeAfter  *string        `gorm:"column:latitude_longitude_after;default:null" json:"latlongAfter"`
	Note                    *string        `gorm:"column:note;default:null" json:"note"`
	IsApprove               *int16         `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt               *string        `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID               *int64         `gorm:"column:approve_id;default:null" json:"approve_id"`
	CreatedAt               string         `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt               string         `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	UserID                  int32          `gorm:"column:user_id;not null" json:"userId"`
}

// TableName CustomerRelocation's table name
func (*CustomerRelocation) TableName() string {
	return TableNameCustomerRelocation
}
