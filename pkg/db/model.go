package db

type CommonModel struct {
	ID        int64 `gorm:"column:id;primarykey"`
	CreatedAt int64 `gorm:"column:created_at"`
	UpdatedAt int64 `gorm:"column:updated_at"`
}
