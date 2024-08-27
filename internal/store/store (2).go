package store

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/models"
	"context"
)

type (
	Store interface {
		Connect(cfg *config.StorageConfig) error
		Close() error
		Title() TitleRepository
		Category() CategoryRepository
		Publisher() PublisherRepository
		User() UserRepository
	}
	TitleRepository interface {
		Create(ctx context.Context, Title *models.Title) error
		All(ctx context.Context, filter *models.Titlesfilter) ([]*models.Title, error)
		ByID(ctx context.Context, id int) (*models.Title, error)
		Update(ctx context.Context, Title *models.Title) error
		Delete(ctx context.Context, id int) error
		ByCategoryId(ctx context.Context, category_id int) (*models.Title, error)
	}
	CategoryRepository interface {
		Create(ctx context.Context, Title *models.Category) error
		Get(ctx context.Context, filter *models.Categoryesfilter) ([]*models.Category, error)
		Update(ctx context.Context, category *models.Category) error
		Delete(ctx context.Context, id int) error
	}
	UserRepository interface {
		Create(ctx context.Context, User *models.User) error
		Update(ctx context.Context, User *models.User) error
		Get(ctx context.Context, filter *models.UserFilter) ([]*models.User, error)
		ByID(ctx context.Context, id int) (*models.User, error)
		Delete(ctx context.Context, id int) error
	}
	PublisherRepository interface {
		Create(ctx context.Context, Publisher *models.Publisher) error
		Get(ctx context.Context, filter *models.Publisherfilter) ([]*models.Publisher, error)
		Update(ctx context.Context, Publisher *models.Publisher) error
		ByID(ctx context.Context, id int) (*models.Publisher, error)
		Delete(ctx context.Context, id int) error
	}
)
