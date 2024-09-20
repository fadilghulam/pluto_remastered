// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package structs

import (
	"time"
)

const TableNameCustomerTypeRequest = "customer_type_request"

// CustomerTypeRequest mapped from table <customer_type_request>
type CustomerTypeRequest struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerTypeID int16     `gorm:"column:customer_type_id;not null" json:"customerTypeId"`
	CustomerID     string    `gorm:"column:customer_id;not null" json:"customerIds"`
	DateEffective  time.Time `gorm:"column:date_effective;not null;default:now()" json:"accessDate"`
	RequestAt      time.Time `gorm:"column:request_at;not null;default:now()" json:"request_at"`
	RequestedID    int64     `gorm:"column:requested_id;default:null" json:"employeeId"`
	IsApprove      int16     `gorm:"column:is_approve;default:null" json:"is_approve"`
	ApproveAt      time.Time `gorm:"column:approve_at;default:null" json:"approve_at"`
	ApproveID      int64     `gorm:"column:approve_id;default:null" json:"approve_id"`
	Note           string    `gorm:"column:note;default:null" json:"note"`
	Attachment     string    `gorm:"column:attachment;default:null" json:"attachment"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

// TableName CustomerTypeRequest's table name
func (*CustomerTypeRequest) TableName() string {
	return TableNameCustomerTypeRequest
}
