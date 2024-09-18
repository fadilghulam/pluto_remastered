package structs

import (
	"time"
)

const TableNameHrEmployee = "hr.employee"

type HrEmployee struct {
	ID                        int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	PersonID                  int64      `gorm:"column:person_id;default:null" json:"person_id"`
	Eid                       string     `gorm:"column:eid;default:null" json:"eid"`
	Email                     string     `gorm:"column:email;default:null" json:"email"`
	Phone                     string     `gorm:"column:phone;default:null" json:"phone"`
	WorkLocationID            int16      `gorm:"column:work_location_id;default:null" json:"work_location_id"`
	DepartmentID              int16      `gorm:"column:department_id;default:null" json:"department_id"`
	DivisionID                int16      `gorm:"column:division_id;default:null" json:"division_id"`
	CreatedAt                 time.Time  `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt                 time.Time  `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
	Salary                    float64    `gorm:"column:salary;default:0" json:"salary"`
	UserID                    int32      `gorm:"column:user_id;default:null" json:"user_id"`
	JobTitleID                int32      `gorm:"column:job_title_id;not null" json:"job_title_id"`
	JobLevelID                int32      `gorm:"column:job_level_id;not null" json:"job_level_id"`
	Photo                     string     `gorm:"column:photo;default:null" json:"photo"`
	Status2                   string     `gorm:"column:status2;default:null" json:"status2"`
	JoinDate                  time.Time  `gorm:"column:join_date;default:null" json:"join_date"`
	IsArchieve                int16      `gorm:"column:is_archieve;default:0;not null" json:"is_archieve"`
	HeadshipID                int64      `gorm:"column:headship_id;default:null" json:"headship_id"`
	AnnualLeaveLeft           int16      `gorm:"column:annual_leave_left;default:12" json:"annual_leave_left"`
	LastActivityAt            time.Time  `gorm:"column:last_activity_at;default:null" json:"last_activity_at"`
	LastActive                time.Time  `gorm:"column:last_active;default:null" json:"last_active"`
	UsePushNotification       int16      `gorm:"column:use_push_notification;default:1;not null" json:"use_push_notification"`
	UseWhatsappNotification   int16      `gorm:"column:use_whatsapp_notification;default:1;not null" json:"use_whatsapp_notification"`
	Status                    string     `gorm:"column:status;default:null" json:"status"`
	SkID                      Int32Array `gorm:"column:sk_id;default:null" json:"sk_id"`
	QrPrintNumber             int16      `gorm:"column:qr_print_number;default:0;not null" json:"qr_print_number"`
	AccessTypeEmployee        string     `gorm:"column:access_type_employee;default:null" json:"access_type_employee"`
	DetailsAccessTypeEmployee Int32Array `gorm:"column:details_access_type_employee;default:null" json:"details_access_type_employee"`
	IsShow                    int16      `gorm:"column:is_show;default:1;not null" json:"is_show"`
	IsSuspend                 int16      `gorm:"column:is_suspend;default:0" json:"is_suspend"`
}

func (*HrEmployee) TableName() string {
	return TableNameHrEmployee
}
