package structs

import (
	"time"
)

const TableNameStokUser = "public.stok_user"

// Stok user mapped from table <stok_user>
type StokUser struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserId          int32     `gorm:"column:user_id" json:"userId"`
	UserIdSubtitute int32     `gorm:"column:user_id_subtitute" json:"userIdSubtitute"`
	TanggalStok     time.Time `gorm:"column:tanggal_stok" json:"tanggalStok"`
	IsComplete      int16     `gorm:"column:is_complete;default: 0" json:"is_complete"`
	TanggalSo       time.Time `gorm:"column:tanggal_so" json:"tanggal_so"`
	SoAdminGudangId int16     `gorm:"column:so_admin_gudang_id" json:"so_admin_gudang_id"`
	ConfirmKey      string    `gorm:"column:confirm_key" json:"confirm_key"`
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;default:now()" json:"updated_at"`
}

// TableName Stok User's table name
func (*StokUser) TableName() string {
	return TableNameStokUser
}
