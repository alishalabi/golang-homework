package handler

import (
	"github.com/labstack/echo/v4"
	"golang-starter-pack/router/middleware"
	"golang-starter-pack/utils"
)

func (h *Handler) Register(v1 *echo.Group) {
	// jwtMiddleware := middleware.JWT(utils.JWTSecret)
	// guestAlbums := v1.Group("/albums")
	// guestAlbums.POST("", h.SignUp)
	// guestUsers.POST("/login", h.Login)

	album := v1.Group("/album")
	album.GET("", h.CurrentAlbum)
	album.PUT("", h.UpdateAlbum)

	profiles := v1.Group("/profiles", jwtMiddleware)
	profiles.GET("/:username", h.GetProfile)
	profiles.POST("/:username/follow", h.Follow)
	profiles.DELETE("/:username/follow", h.Unfollow)

	articles := v1.Group("/articles", middleware.JWTWithConfig(
		middleware.JWTConfig{
			Skipper: func(c echo.Context) bool {
				if c.Request().Method == "GET" && c.Path() != "/api/articles/feed" {
					return true
				}
				return false
			},
			SigningKey: utils.JWTSecret,
		},
	))
	articles.POST("", h.CreateArticle)
	articles.GET("/feed", h.Feed)
	articles.PUT("/:slug", h.UpdateArticle)
	articles.DELETE("/:slug", h.DeleteArticle)
	articles.POST("/:slug/comments", h.AddComment)
	articles.DELETE("/:slug/comments/:id", h.DeleteComment)
	articles.POST("/:slug/favorite", h.Favorite)
	articles.DELETE("/:slug/favorite", h.Unfavorite)
	articles.GET("", h.Articles)
	articles.GET("/:slug", h.GetArticle)
	articles.GET("/:slug/comments", h.GetComments)

	tags := v1.Group("/tags")
	tags.GET("", h.Tags)
}
