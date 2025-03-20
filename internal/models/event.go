package models

type Event struct {
	Type              int64                     `json:"type"` // 0 - монолог героя или закадровый голос, 1- персонаж появилсяб 2 - перслнаж ушел, 3 - персонаж произносит речь
	Character         int64                     `json:"character"`
	Sound             int64                     `json:"sound"`
	CharactersInEvent map[int64]map[int64]int64 `json:"characters_in_event" gorm:"type:json"` // персонаж - эмоция + позиция относительно левого края
	Text              string                    `json:"text"`
}
