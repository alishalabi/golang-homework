package handler

import (
	"github.com/labstack/echo/v4"
	"golang-starter-pack/router/middleware"
	"golang-starter-pack/utils"
)

func (h *Handler) Register(v1 *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	guestUsers := v1.Group("/users")
	guestUsers.POST("", h.SignUp)
	guestUsers.POST("/login", h.Login)

	user := v1.Group("/user", jwtMiddleware)
	user.GET("", h.CurrentUser)
	user.PUT("", h.UpdateUser)

	profiles := v1.Group("/profiles", jwtMiddleware)
	profiles.GET("/:username", h.GetProfile)
	profiles.POST("/:username/follow", h.Follow)
	profiles.DELETE("/:username/follow", h.Unfollow)

	projects := v1.Group("/projects", middleware.JWTWithConfig(
		middleware.JWTConfig{
			Skipper: func(c echo.Context) bool {
				if c.Request().Method == "GET" && c.Path() != "/api/projects/feed" {
					return true
				}
				return false
			},
			SigningKey: utils.JWTSecret,
		},
	))
	projects.POST("", h.CreateProject)
	projects.GET("/feed", h.Feed)
	projects.PUT("/:slug", h.UpdateProject)
	projects.DELETE("/:slug", h.DeleteProject)
	projects.POST("/:slug/favorite", h.Favorite)
	projects.DELETE("/:slug/favorite", h.Unfavorite)
	projects.GET("", h.Projects)
	projects.GET("/:slug", h.GetProject)

	tags := v1.Group("/tags")
	tags.GET("", h.Tags)
}
