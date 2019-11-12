package handler

import (
	"net/http"

	"golang-starter-pack/model"
	"golang-starter-pack/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) NewAlbum(c echo.Context) error {
	var u model.Album
	req := &albumRegisterRequest{}
	if err := req.bind(c, &u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := h.albumStore.Create(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, newAlbumResponse(&u))
}


func (h *Handler) UpdateAlbum(c echo.Context) error {
	u, err := h.albumStore.GetByID(albumIDFromToken(c))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	req := newAlbumUpdateRequest()
	req.populate(u)
	if err := req.bind(c, u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := h.albumStore.Update(u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newAlbumResponse(u))
}

func (h *Handler) GetProfile(c echo.Context) error {
	title := c.Param("title")
	u, err := h.albumStore.GetByTitle(title)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newProfileResponse(h.albumStore, albumIDFromToken(c), u))
}
