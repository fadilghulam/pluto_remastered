package structs

const TableNameUserSession = "public.user_session"

// User Session mapped from table <public.user_session>
type UserSession struct {
	ID               int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID           *int32 `gorm:"column:user_id" json:"user_id"`
	AppVersion       string `gorm:"column:app_version;default: null" json:"app_version"`
	BatteryLevel     string `gorm:"column:battery_level;default: null" json:"battery_level"`
	CarrierName      string `gorm:"column:carrier_name;default: null" json:"carrier_name"`
	DeviceBrand      string `gorm:"column:device_brand;default: null" json:"device_brand"`
	DeviceID         string `gorm:"column:device_id;default: null" json:"device_id"`
	DeviceName       string `gorm:"column:device_name;default: null" json:"device_name"`
	IP               string `gorm:"column:ip;default: null" json:"ip"`
	IsMock           *int16 `gorm:"column:is_mock;default: null" json:"is_mock"`
	Latitude         string `gorm:"column:latitude;default: null" json:"latitude"`
	Longitude        string `gorm:"column:longitude;default: null" json:"longitude"`
	Language         string `gorm:"column:language;default: null" json:"language"`
	LoginMethod      string `gorm:"column:login_method;default: null" json:"login_method"`
	NetworkType      string `gorm:"column:network_type;default: null" json:"network_type"`
	OsName           string `gorm:"column:os_name;default: null" json:"os_name"`
	OsVersion        string `gorm:"column:os_version;default: null" json:"os_version"`
	Platform         string `gorm:"column:platform;default: null" json:"platform"`
	ScreenResolution string `gorm:"column:screen_resolution;default: null" json:"screen_resolution"`
	Timestamp        string `gorm:"column:timestamp;default: null" json:"timestamp"`
	Timezone         string `gorm:"column:timezone;default: null" json:"timezone"`
	UserIDSubtitute  int32  `gorm:"column:user_id_subtitute;default: null" json:"user_id_subtitute"`
}

// TableName Salesman's table name
func (*UserSession) TableName() string {
	return TableNameUserSession
}
