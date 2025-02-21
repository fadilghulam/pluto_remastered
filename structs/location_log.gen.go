package structs

const TableNameLocationLog = "location_log"

// LocationLog mapped from table <location_log>
type LocationLog struct {
	ID                              FlexibleString `gorm:"column:id;primaryKey" json:"id"`
	UserID                          *int32         `gorm:"column:user_id;default:null" json:"user_id"`
	Title                           *string        `gorm:"column:title;default:null;comment:ex: KUNJUNGAN, PENJUALAN" json:"title"` // ex: KUNJUNGAN, PENJUALAN
	Subtitle                        *string        `gorm:"column:subtitle;default:null" json:"subtitle"`
	SalesmanID                      *int32         `gorm:"column:salesman_id;default:null" json:"salesman_id"`
	CustomerID                      FlexibleString `gorm:"column:customer_id;default:null" json:"customer_id"`
	Latitude                        *float64       `gorm:"column:latitude;default:null" json:"latitude"`
	Longitude                       *float64       `gorm:"column:longitude;default:null" json:"longitude"`
	Accuracy                        *float64       `gorm:"column:accuracy;default:null" json:"accuracy"`
	VerticalAccuracy                *float64       `gorm:"column:vertical_accuracy;default:null" json:"vertical_accuracy"`
	Altitude                        *float64       `gorm:"column:altitude;default:null" json:"altitude"`
	Speed                           *float64       `gorm:"column:speed;default:null" json:"speed"`
	SpeedAccuracy                   *float64       `gorm:"column:speed_accuracy;default:null" json:"speed_accuracy"`
	Heading                         *float64       `gorm:"column:heading;default:null" json:"heading"`
	Time                            *float64       `gorm:"column:time;default:null" json:"time"`
	IsMock                          *int16         `gorm:"column:is_mock;default:null" json:"is_mock"`
	HeadingAccuracy                 *float64       `gorm:"column:heading_accuracy;default:null" json:"heading_accuracy"`
	ElapsedRealtimeNanos            *float64       `gorm:"column:elapsed_realtime_nanos;default:null" json:"elapsed_realtime_nanos"`
	ElapsedRealtimeUncertaintyNanos *float64       `gorm:"column:elapsed_realtime_uncertainty_nanos;default:null" json:"elapsed_realtime_uncertainty_nanos"`
	SatelliteNumber                 *int32         `gorm:"column:satellite_number;default:null" json:"satellite_number"`
	Provider                        *string        `gorm:"column:provider;default:null" json:"provider"`
	CreatedAt                       string         `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt                       string         `gorm:"column:updated_at;default:now()" json:"updated_at"`
	SyncKey                         string         `gorm:"column:sync_key;default:now()" json:"sync_key"`
	Label                           *string        `gorm:"column:label;default:null" json:"label"`
}

// TableName LocationLog's table name
func (*LocationLog) TableName() string {
	return TableNameLocationLog
}
