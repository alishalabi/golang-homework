package handler

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"golang-starter-pack/model"
)

type userUpdateRequest struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email" validate:"email"`
		Password string `json:"password"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user"`
}

func newUserUpdateRequest() *userUpdateRequest {
	return new(userUpdateRequest)
}

func (r *userUpdateRequest) populate(u *model.User) {
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Password = u.Password
	if u.Bio != nil {
		r.User.Bio = *u.Bio
	}
	if u.Image != nil {
		r.User.Image = *u.Image
	}
}

func (r *userUpdateRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.User.Username
	u.Email = r.User.Email
	if r.User.Password != u.Password {
		h, err := u.HashPassword(r.User.Password)
		if err != nil {
			return err
		}
		u.Password = h
	}
	u.Bio = &r.User.Bio
	u.Image = &r.User.Image
	return nil
}

type userRegisterRequest struct {
	User struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userRegisterRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.User.Username
	u.Email = r.User.Email
	h, err := u.HashPassword(r.User.Password)
	if err != nil {
		return err
	}
	u.Password = h
	return nil
}

type userLoginRequest struct {
	User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userLoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type projectCreateRequest struct {
	Project struct {
		Title       string   `json:"title" validate:"required"`
		Description string   `json:"description" validate:"required"`
		Body        string   `json:"body" validate:"required"`
		Tags        []string `json:"tagList, omitempty"`
	} `json:"project"`
}

func (r *projectCreateRequest) bind(c echo.Context, a *model.Project) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	a.Title = r.Project.Title
	a.Slug = slug.Make(r.Project.Title)
	a.Description = r.Project.Description
	a.Body = r.Project.Body
	if r.Project.Tags != nil {
		for _, t := range r.Project.Tags {
			a.Tags = append(a.Tags, model.Tag{Tag: t})
		}
	}
	return nil
}

type projectUpdateRequest struct {
	Project struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Body        string   `json:"body"`
		Tags        []string `json:"tagList"`
	} `json:"project"`
}

func (r *projectUpdateRequest) populate(a *model.Project) {
	r.Project.Title = a.Title
	r.Project.Description = a.Description
	r.Project.Body = a.Body
}

func (r *projectUpdateRequest) bind(c echo.Context, a *model.Project) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	a.Title = r.Project.Title
	a.Slug = slug.Make(a.Title)
	a.Description = r.Project.Description
	a.Body = r.Project.Body
	return nil
}
