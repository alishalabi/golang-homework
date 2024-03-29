package model

import (
	"github.com/jinzhu/gorm"
)

type Project struct {
	gorm.Model
	Slug        string `gorm:"unique_index;not null"`
	Title       string `gorm:"not null"`
	Description string
	Body        string
	Author      string
	// AuthorID    uint
	// Tags        []Tag  `gorm:"many2many:project_tags;association_autocreate:false"`
}


type Tag struct {
	gorm.Model
	Tag      string    `gorm:"unique_index"`
	Projects []Project `gorm:"many2many:project_tags;"`
}
