package domain

import "time"

type Player struct {
	ID           int64
	Name         string
	IconPath     string
	CoverImgPath string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
