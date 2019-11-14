package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang-starter-pack/model"
	"golang-starter-pack/utils"
)

func (h *Handler) GetProject(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.projectStore.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newProjectResponse(c, a))
}

func (h *Handler) Projects(c echo.Context) error {
	tag := c.QueryParam("tag")
	author := c.QueryParam("author")
	favoritedBy := c.QueryParam("favorited")
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 20
	}
	var projects []model.Project
	var count int
	if tag != "" {
		projects, count, err = h.projectStore.ListByTag(tag, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else if author != "" {
		projects, count, err = h.projectStore.ListByAuthor(author, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else if favoritedBy != "" {
		projects, count, err = h.projectStore.ListByWhoFavorited(favoritedBy, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else {
		projects, count, err = h.projectStore.List(offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return c.JSON(http.StatusOK, newProjectListResponse(h.userStore, userIDFromToken(c), projects, count))
}

func (h *Handler) Feed(c echo.Context) error {
	var projects []model.Project
	var count int
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 20
	}
	projects, count, err = h.projectStore.ListFeed(userIDFromToken(c), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, newProjectListResponse(h.userStore, userIDFromToken(c), projects, count))
}

func (h *Handler) CreateProject(c echo.Context) error {
	var a model.Project
	req := &projectCreateRequest{}
	if err := req.bind(c, &a); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	a.AuthorID = userIDFromToken(c)
	err := h.projectStore.CreateProject(&a)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, newProjectResponse(c, &a))
}

func (h *Handler) UpdateProject(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.projectStore.GetUserProjectBySlug(userIDFromToken(c), slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	req := &projectUpdateRequest{}
	req.populate(a)
	if err := req.bind(c, a); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err = h.projectStore.UpdateProject(a, req.Project.Tags); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newProjectResponse(c, a))
}

func (h *Handler) DeleteProject(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.projectStore.GetUserProjectBySlug(userIDFromToken(c), slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	err = h.projectStore.DeleteProject(a)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

func (h *Handler) AddComment(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.projectStore.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	var cm model.Comment
	req := &createCommentRequest{}
	if err := req.bind(c, &cm); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err = h.projectStore.AddComment(a, &cm); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, newCommentResponse(c, &cm))
}

func (h *Handler) GetComments(c echo.Context) error {
	slug := c.Param("slug")
	cm, err := h.projectStore.GetCommentsBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newCommentListResponse(c, cm))
}

func (h *Handler) DeleteComment(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}
	cm, err := h.projectStore.GetCommentByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if cm == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if cm.UserID != userIDFromToken(c) {
		return c.JSON(http.StatusUnauthorized, utils.NewError(errors.New("unauthorized action")))
	}
	if err := h.projectStore.DeleteComment(cm); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

func (h *Handler) Favorite(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.projectStore.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if err := h.projectStore.AddFavorite(a, userIDFromToken(c)); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newProjectResponse(c, a))
}

func (h *Handler) Unfavorite(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.projectStore.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if err := h.projectStore.RemoveFavorite(a, userIDFromToken(c)); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newProjectResponse(c, a))
}

func (h *Handler) Tags(c echo.Context) error {
	tags, err := h.projectStore.ListTags()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, newTagListResponse(tags))
}
