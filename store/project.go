package store

import (
	"github.com/jinzhu/gorm"
	"golang-starter-pack/model"
)

type ProjectStore struct {
	db *gorm.DB
}

func NewProjectStore(db *gorm.DB) *ProjectStore {
	return &ProjectStore{
		db: db,
	}
}

func (ps *ProjectStore) GetBySlug(s string) (*model.Project, error) {
	var m model.Project
	err := ps.db.Where(&model.Project{Slug: s}).Preload("Favorites").Preload("Tags").Preload("Author").Find(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (ps *ProjectStore) GetUserProjectBySlug(userID uint, slug string) (*model.Project, error) {
	var m model.Project
	err := ps.db.Where(&model.Project{Slug: slug, AuthorID: userID}).Find(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (ps *ProjectStore) CreateProject(a *model.Project) error {
	tags := a.Tags
	tx := ps.db.Begin()
	if err := tx.Create(&a).Error; err != nil {
		return err
	}
	for _, t := range a.Tags {
		err := tx.Where(&model.Tag{Tag: t.Tag}).First(&t).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&a).Association("Tags").Append(t).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where(a.ID).Preload("Favorites").Preload("Tags").Preload("Author").Find(&a).Error; err != nil {
		tx.Rollback()
		return err
	}
	a.Tags = tags
	return tx.Commit().Error
}

func (ps *ProjectStore) UpdateProject(a *model.Project, tagList []string) error {
	tx := ps.db.Begin()
	if err := tx.Model(a).Update(a).Error; err != nil {
		return err
	}
	tags := make([]model.Tag, 0)
	for _, t := range tagList {
		tag := model.Tag{Tag: t}
		err := tx.Where(&tag).First(&tag).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			return err
		}
		tags = append(tags, tag)
	}
	if err := tx.Model(a).Association("Tags").Replace(tags).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(a.ID).Preload("Favorites").Preload("Tags").Preload("Author").Find(a).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (ps *ProjectStore) DeleteProject(a *model.Project) error {
	return ps.db.Delete(a).Error
}

func (ps *ProjectStore) List(offset, limit int) ([]model.Project, int, error) {
	var (
		projects []model.Project
		count    int
	)
	ps.db.Model(&projects).Count(&count)
	ps.db.Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Find(&projects)
	return projects, count, nil
}

func (ps *ProjectStore) ListByTag(tag string, offset, limit int) ([]model.Project, int, error) {
	var (
		t        model.Tag
		projects []model.Project
		count    int
	)
	err := ps.db.Where(&model.Tag{Tag: tag}).First(&t).Error
	if err != nil {
		return nil, 0, err
	}
	ps.db.Model(&t).Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Association("Projects").Find(&projects)
	count = ps.db.Model(&t).Association("Projects").Count()
	return projects, count, nil
}

func (ps *ProjectStore) ListByAuthor(username string, offset, limit int) ([]model.Project, int, error) {
	var (
		u        model.User
		projects []model.Project
		count    int
	)
	err := ps.db.Where(&model.User{Username: username}).First(&u).Error
	if err != nil {
		return nil, 0, err
	}
	ps.db.Where(&model.Project{AuthorID: u.ID}).Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Find(&projects)
	ps.db.Where(&model.Project{AuthorID: u.ID}).Model(&model.Project{}).Count(&count)

	return projects, count, nil
}

func (ps *ProjectStore) ListByWhoFavorited(username string, offset, limit int) ([]model.Project, int, error) {
	var (
		u        model.User
		projects []model.Project
		count    int
	)
	err := ps.db.Where(&model.User{Username: username}).First(&u).Error
	if err != nil {
		return nil, 0, err
	}
	ps.db.Model(&u).Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Association("Favorites").Find(&projects)
	count = ps.db.Model(&u).Association("Favorites").Count()
	return projects, count, nil
}

func (ps *ProjectStore) ListFeed(userID uint, offset, limit int) ([]model.Project, int, error) {
	var (
		u        model.User
		projects []model.Project
		count    int
	)
	err := ps.db.First(&u, userID).Error
	if err != nil {
		return nil, 0, err
	}
	var followings []model.Follow
	ps.db.Model(&u).Preload("Following").Preload("Follower").Association("Followings").Find(&followings)
	var ids []uint
	for _, i := range followings {
		ids = append(ids, i.FollowingID)
	}
	ps.db.Where("author_id in (?)", ids).Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Find(&projects)
	ps.db.Where(&model.Project{AuthorID: u.ID}).Model(&model.Project{}).Count(&count)
	return projects, count, nil
}


func (ps *ProjectStore) AddFavorite(a *model.Project, userID uint) error {
	usr := model.User{}
	usr.ID = userID
	return ps.db.Model(a).Association("Favorites").Append(&usr).Error
}

func (ps *ProjectStore) RemoveFavorite(a *model.Project, userID uint) error {
	usr := model.User{}
	usr.ID = userID
	return ps.db.Model(a).Association("Favorites").Delete(&usr).Error
}

func (ps *ProjectStore) ListTags() ([]model.Tag, error) {
	var tags []model.Tag
	if err := ps.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
