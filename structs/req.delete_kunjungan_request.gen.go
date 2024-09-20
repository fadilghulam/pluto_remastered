// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package structs

import (
	"time"
)

const TableNameDeleteKunjunganRequest = "delete_kunjungan_request"

// DeleteKunjunganRequest mapped from table <delete_kunjungan_request>
type DeleteKunjunganRequest struct {
	EmployeeID  int64     `gorm:"column:employee_id;not null" json:"employeeId"`
	Datetime    time.Time `gorm:"column:datetime;not null;default:now()" json:"datetime"`
	KunjunganID int64     `gorm:"column:kunjungan_id;not null" json:"kunjunganId"`
	Note        string    `gorm:"column:note;default:null" json:"note"`
	IsApprove   int16     `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt   time.Time `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID   int64     `gorm:"column:approve_id;default:null" json:"approve_id"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID 		int32 	  `gorm:"column:user_id;not null" json:"userId"`
}

// TableName DeleteKunjunganRequest's table name
func (*DeleteKunjunganRequest) TableName() string {
	return TableNameDeleteKunjunganRequest
}
