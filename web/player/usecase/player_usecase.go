package usecase

import (
	"context"
	"time"

	"github.com/itouri/fortnite/web/domain"
	"github.com/itouri/fortnite/web/player"
)

type playerUsecase struct {
	playerRepo     player.Repository
	contextTimeout time.Duration
}

func NewPlayerUsecase(p player.Repository, timeout time.Duration) playerUsecase {
	return &playerUsecase{
		playerRepo:     p,
		contextTimeout: timeout,
	}
}

// fillAuthorDetails?

func (p *playerUsecase) Fetch(ctx context.Context, cursor string, num int64) ([]*domain.Player, string, error) {
	//LEARN
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	listPlayer, nextCursor, err := p.playerRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	return listPlayer, nextCursor, nil
}

func (p *playerUsecase) GetByID(ctx context.Context, id int64) (*domain.Player, error) {

	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	res, err := p.playerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *playerUsecase) Update(ctx context.Context, pl *domain.Player) error {

	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	pl.UpdatedAt = time.Now()
	return p.playerRepo.Update(ctx, pl)
}

func (p *playerUsecase) Store(ctx context.Context, pl *domain.Player) error {

	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

}
