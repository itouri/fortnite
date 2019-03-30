package domain

import "time"

type Recruit struct {
	ID                int64
	IsPublic          bool
	Atomsphere        int64 //TINYINT なのに int64 は過剰
	RecruitmentNumber int64 //TINYINT
	BeginTime         time.Time
	EndTime           time.Time
	Message           string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
