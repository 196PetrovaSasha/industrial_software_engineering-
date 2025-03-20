package models

type Player struct {
	Id                int64 `gorm:"primary_key"`
	Name              string
	Email             string
	Phone             string
	Password          string
	Admin             bool
	CompletedChapters []int64         `gorm:"type:json;column:completed_chapters"` // пройденные главы
	ChaptersProgress  map[int64]int64 `gorm:"type:json;column:chapters_progress"`  // Мапа id главы - id узла
	SoundSettings     int             //percent
}
