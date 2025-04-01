package storage

import (
	"context"
	"gopkg.in/reform.v1"
)

type (
	Storage interface {
		CreateNew(ctx context.Context, new New) error
		GetNews(ctx context.Context) ([]New, error)
		GetNewFromID(ctx context.Context, id int) (New, error)
		UpdateNewFromID(ctx context.Context, oldID int, new New) error
	}
	PostgresStorage struct {
		db *reform.DB
	}
	New struct {
		Id         int    `json:"Id"`
		Title      string `json:"Title"`
		Content    string `json:"Content"`
		Categories []int  `json:"Categories"`
	}
)
