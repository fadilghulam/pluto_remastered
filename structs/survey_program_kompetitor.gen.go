package structs

const TableNameSurveyProgramKompetitor = "survey_program_kompetitor"

// SurveyProgramKompetitor mapped from table <survey_program_kompetitor>
type SurveyProgramKompetitor struct {
	ID                FlexibleString `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CustomerID        FlexibleString `gorm:"column:customer_id;default:null" json:"customer_id"`
	SalesmanID        *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	MerchandiserID    *int32         `gorm:"column:merchandiser_id;default:null" json:"merchandiser_id"`
	TeamleaderID      *int32         `gorm:"column:teamleader_id;default:null" json:"teamleader_id"`
	Datetime          *string        `gorm:"column:datetime;default:null" json:"datetime"`
	ProgramName       *string        `gorm:"column:program_name;default:null" json:"program_name"`
	Period            *string        `gorm:"column:period;default:null" json:"period"`
	StartDate         *string        `gorm:"column:start_date;default:null" json:"start_date"`
	EndDate           *string        `gorm:"column:end_date;default:null" json:"end_date"`
	Information       *string        `gorm:"column:information;default:null" json:"information"`
	CreatedAt         string         `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt         string         `gorm:"column:updated_at;default:now()" json:"updated_at"`
	SyncKey           string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	Photo             *string        `gorm:"column:photo;default:null" json:"photo"`
	LatitudeLongitude *string        `gorm:"column:latitude_longitude;default:null" json:"latitude_longitude"`
	KompetitorID      FlexibleString `gorm:"column:kompetitor_id;default:null" json:"kompetitor_id"`
}

// TableName SurveyProgramKompetitor's table name
func (*SurveyProgramKompetitor) TableName() string {
	return TableNameSurveyProgramKompetitor
}
