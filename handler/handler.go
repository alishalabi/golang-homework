package handler

import (
	"golang-starter-pack/project"
	"golang-starter-pack/user"
)

type Handler struct {
	userStore    user.Store
	projectStore project.Store
}

func NewHandler(us user.Store, ps project.Store) *Handler {
	return &Handler{
		userStore:    us,
		projectStore: ps,
	}
}
