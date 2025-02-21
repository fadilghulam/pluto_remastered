package structs

const TableNameAppVersioning = "app_versioning"

// AppVersioning mapped from table <app_versioning>
type AppVersioning struct {
	ID                    int64   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	AppName               *string `gorm:"column:app_name;default:null;comment:version" json:"app_name"` // version
	CurrentVersion        *string `gorm:"column:current_version;default:null" json:"current_version"`
	MinimumVersion        *string `gorm:"column:minimum_version;default:null" json:"minimum_version"`
	ForceUpdate           *int16  `gorm:"column:force_update;default:null" json:"force_update"`
	Changelog             *string `gorm:"column:changelog;default:null" json:"changelog"`
	AndroidURL            *string `gorm:"column:android_url;default:null" json:"android_url"`
	IosURL                *string `gorm:"column:ios_url;default:null" json:"ios_url"`
	IsMaintenance         *int16  `gorm:"column:is_maintenance;default:null;comment:maintenance" json:"is_maintenance"` // maintenance
	Message               *string `gorm:"column:message;default:null" json:"message"`
	MinimumAndroidVersion *string `gorm:"column:minimum_android_version;default:null;comment:os_compability" json:"minimum_android_version"` // os_compability
	MinimumIosVersion     *string `gorm:"column:minimum_ios_version;default:null" json:"minimum_ios_version"`
	APIBaseURL            *string `gorm:"column:api_base_url;default:null;comment:configuration" json:"api_base_url"` // configuration
	RequestTimeout        *int16  `gorm:"column:request_timeout;default:null" json:"request_timeout"`
	MaxUploadSizeMb       *int16  `gorm:"column:max_upload_size_mb;default:null" json:"max_upload_size_mb"`
	TermsAndConditions    *string `gorm:"column:terms_and_conditions;default:null;comment:urls" json:"terms_and_conditions"` // urls
	PrivacyPolicy         *string `gorm:"column:privacy_policy;default:null" json:"privacy_policy"`
	Landing               *string `gorm:"column:landing;default:null" json:"landing"`
}

// TableName AppVersioning's table name
func (*AppVersioning) TableName() string {
	return TableNameAppVersioning
}
