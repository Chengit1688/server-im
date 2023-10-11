package model

type FavoriteImage struct {
	ID             int64  `gorm:"column:id;primarykey"`
	UserID         string `gorm:"column:user_id;size:80;index;" json:"user_id"`
	UUID           string `gorm:"column:uuid;size:64;index;uniqueIndex:idx_user_id_uuid;" json:"uuid"`
	ImageUrl       string `gorm:"column:image_url;size:240" json:"image_url"`
	ImageThumbnail string `gorm:"column:image_thumbnail;size:240" json:"image_thumbnail"`
	ImageWidth     *int   `gorm:"column:image_width" json:"image_width"`
	ImageHeight    *int   `gorm:"column:image_height" json:"image_height"`
}

func (t *FavoriteImage) TableName() string {
	return "favorite_image"
}
