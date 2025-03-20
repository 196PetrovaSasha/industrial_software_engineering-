package models

type Media struct {
	Id          int64  `gorm:"primarykey"`
	FileData    []byte `json:"-" gorm:"type:bytea;column:file_data"` // для хранения файла
	ContentType string `json:"content_type"`
}
