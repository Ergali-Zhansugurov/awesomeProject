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
		All(ctx context.Context, filter *models.Filter) ([]*models.Title, error)
		ByID(ctx context.Context, id int) (*models.Title, error)
		Update(ctx context.Context, Title *models.Title) error
		Delete(ctx context.Context, id int) error
		ByCategoryId(ctx context.Context, category_id int) (*models.Title, error)
	}
	CategoryRepository interface {
		Create(ctx context.Context, Title *models.Category) error
		Get(ctx context.Context, filter *models.Filter) ([]*models.Category, error)
		Update(ctx context.Context, category *models.Category) error
		Delete(ctx context.Context, id int) error
	}
	UserRepository interface {
		Create(ctx context.Context, User *models.User) error
		Update(ctx context.Context, User *models.User) error
		Get(ctx context.Context, filter *models.Filter) ([]*models.User, error)
		ByID(ctx context.Context, id int) (*models.User, error)
		Delete(ctx context.Context, id int) error
		AddTitleToLibrary(ctx context.Context, userID int, title *models.Title) error
		RemoveTitleFromLibrary(ctx context.Context, userID int, titleID int) error
		Authenticate(ctx context.Context, username, password string) (*models.User, error)
	}
	PublisherRepository interface {
		Create(ctx context.Context, Publisher *models.Publisher) error
		Get(ctx context.Context, filter *models.Filter) ([]*models.Publisher, error)
		Update(ctx context.Context, Publisher *models.Publisher) error
		ByID(ctx context.Context, id int) (*models.Publisher, error)
		Delete(ctx context.Context, id int) error
		AddTitleToLibrary(ctx context.Context, publisherID int, title *models.Title) error
		RemoveTitleFromLibrary(ctx context.Context, publisherID int, titleID int) error
		Authenticate(ctx context.Context, publishername, password string) (*models.User, error)
	}
)
