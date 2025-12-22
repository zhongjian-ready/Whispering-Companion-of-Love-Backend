package model

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID           int64          `gorm:"primaryKey;autoIncrement;type:bigserial;column:id" json:"user_id"`
	Username     *string        `gorm:"type:varchar(50);unique" json:"username,omitempty"`
	Email        *string        `gorm:"type:varchar(100);unique;index:idx_users_email" json:"email,omitempty"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Salt         string         `gorm:"type:varchar(50)" json:"-"`
	OpenID       *string        `gorm:"column:openid;type:varchar(100);unique;index:idx_users_openid" json:"openid,omitempty"`
	UnionID      *string        `gorm:"column:unionid;type:varchar(100)" json:"unionid,omitempty"`
	WechatInfo   datatypes.JSON `gorm:"type:jsonb" json:"wechat_info,omitempty"`
	Nickname     string         `gorm:"type:varchar(100)" json:"nickname"`
	AvatarURL    string         `gorm:"column:avatar_url;type:varchar(500)" json:"avatar_url"`
	Phone        *string        `gorm:"type:varchar(20);index:idx_users_phone" json:"phone,omitempty"`
	Gender       int8           `gorm:"type:smallint;default:0" json:"gender"`
	IsActive     bool           `gorm:"default:true;index:idx_users_is_active" json:"is_active"`
	IsAdmin      bool           `gorm:"default:false" json:"is_admin"`
	IsDeleted    bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt    time.Time      `gorm:"default:CURRENT_TIMESTAMP;index:idx_users_created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	LastLoginIP  string         `gorm:"column:last_login_ip;type:inet" json:"last_login_ip"`
	LoginCount   int            `gorm:"default:0" json:"login_count"`
	Source       string         `gorm:"type:varchar(20);default:'wechat_mini'" json:"source"`

	// Settings
	DailyGoal         int64          `gorm:"default:2000" json:"daily_goal"`
	ReminderSettings  datatypes.JSON `gorm:"type:jsonb" json:"reminder_settings,omitempty"`
	QuickAddSettings  datatypes.JSON `gorm:"type:jsonb" json:"quick_add_settings,omitempty"`
	ReminderEnabled   bool           `gorm:"default:false" json:"reminder_enabled"`
	ReminderInterval  int64          `gorm:"default:60" json:"reminder_interval"`                        // minutes
	ReminderStartTime string         `gorm:"type:varchar(5);default:'08:00'" json:"reminder_start_time"` // "08:00"
	ReminderEndTime   string         `gorm:"type:varchar(5);default:'22:00'" json:"reminder_end_time"`   // "22:00"
	QuickAddPresets   datatypes.JSON `gorm:"type:jsonb" json:"quick_add_presets"`                        // [200, 300, 500, 800]
}

func (User) TableName() string {
	return "users"
}
