package project

import (
	"golang-starter-pack/model"
)

type Store interface {
	GetBySlug(string) (*model.Project, error)
	GetUserProjectBySlug(userID uint, slug string) (*model.Project, error)
	CreateProject(*model.Project) error
	UpdateProject(*model.Project, []string) error
	DeleteProject(*model.Project) error
	List(offset, limit int) ([]model.Project, int, error)
	ListByTag(tag string, offset, limit int) ([]model.Project, int, error)
	ListByAuthor(username string, offset, limit int) ([]model.Project, int, error)
	ListByWhoFavorited(username string, offset, limit int) ([]model.Project, int, error)
	ListFeed(userID uint, offset, limit int) ([]model.Project, int, error)

	AddFavorite(*model.Project, uint) error
	RemoveFavorite(*model.Project, uint) error
	ListTags() ([]model.Tag, error)
}
