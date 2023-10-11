package model

type DashboardDailyData struct {
	ID       int64 `gorm:"column:id;primaryKey"`
	DateTime int64 `gorm:"column:datetime;Index:datetime"`
	Count    int64 `gorm:"column:count"`
	Type     int64 `gorm:"column:type;Index:type;comment:1私聊消息数 2群聊消息数 3在线用户数 4注册用户数"`
}

func (DashboardDailyData) TableName() string {
	return "dashboard_daily_data"
}
