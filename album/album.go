package album

import (
	"golang-starter-pack/model"
)

type Store interface {
	GetByID(uint) (*model.Album, error)
	GetByTitle(string) (*model.Album, error)
	Create(*model.Album) error
	Update(*model.Album) error
}
