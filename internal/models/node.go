package models

type Node struct {
	Id         int64 `gorm:"primary_key"`
	Slug       string
	Events     map[int]Event `gorm:"type:json"`
	ChapterId  int64
	Music      int64
	Background int64
	Branching  Branching `gorm:"type:json"`
	End        EndInfo   `gorm:"type:json"`
	Comment    string
}

type Branching struct {
	Flag      bool
	Condition map[string]int64 `gorm:"type:json"`
}

type EndInfo struct {
	Flag      bool
	EndResult string
	EndText   string
}
