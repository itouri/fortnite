package player

import (
	"context"

	"github.com/itouri/fortnite/web/domain"
)

type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]*domain.Player, string, error)
	GetByID(ctx context.Context, id int64) (*domain.Player, error)
	Update(ctx context.Context, player *domain.Player) error
	Store(ctx context.Context, player *domain.Player) error
	Delete(ctx context.Context, id int64) error
}
