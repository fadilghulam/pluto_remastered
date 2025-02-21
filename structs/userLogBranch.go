package structs

const TableNameUserLogBranch = "public.user_log_branch"

// User Log Branch mapped from table <public.user_log_branch>
type UserLogBranch struct {
	ID              int64   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID          *int32  `gorm:"column:user_id" json:"user_id"`
	UserIDSubtitute *int32  `gorm:"column:user_id_subtitute" json:"user_id_subtitute"`
	BranchID        *int16  `gorm:"column:branch_id" json:"branch_id"`
	StartDate       *string `gorm:"column:start_date" json:"start_date"`
	EndDate         *string `gorm:"column:end_date" json:"end_date"`
	LastVisitDate   *string `gorm:"column:last_visit_date" json:"last_visit_date"`
}

// TableName Salesman's table name
func (*UserLogBranch) TableName() string {
	return TableNameUserLogBranch
}
