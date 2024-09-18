package structs

import (
	"time"
)

const TableNameHrApprovalHistory = "hr.approval_history"

type HrApprovalHistory struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Datetime       time.Time `gorm:"column:datetime;not null;default:now()" json:"datetime"`
	EmployeeID     int64     `gorm:"column:employee_id;not null" json:"employee_id"`
	ReferenceTable string    `gorm:"column:reference_table; not null" json:"reference_table"`
	ReferenceID    int64     `gorm:"column:reference_id; not null" json:"reference_id"`
	Note           string    `gorm:"column:note;default:null" json:"note"`
	IsApprove      int16     `gorm:"column:is_approve;default:0;not null" json:"is_approve"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	Index          int64     `gorm:"column:index;not null" json:"index"`
}

func (*HrApprovalHistory) TableName() string {
	return TableNameHrApprovalHistory
}
