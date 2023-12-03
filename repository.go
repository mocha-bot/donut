package main

import (
	"context"

	"gorm.io/gorm"
)

type donutRepository struct {
	db *gorm.DB
}

type DonutRepository interface {
	CreateMatchMaker(ctx context.Context, matchMaker *MatchMakerEntity) error
	CreateMatchMakerUser(ctx context.Context, matchMakerUser *MatchMakerUser) error
	GetMatchMakerBySerial(ctx context.Context, serial string) (*MatchMakerEntity, error)
}

func NewDonutRepository(db *gorm.DB) DonutRepository {
	return &donutRepository{
		db: db,
	}
}

func (r *donutRepository) CreateMatchMaker(ctx context.Context, matchMaker *MatchMakerEntity) error {
	return r.db.WithContext(ctx).Create(MatchMaker{}.FromEntity(matchMaker)).Error
}

func (r *donutRepository) CreateMatchMakerUser(ctx context.Context, matchMakerUser *MatchMakerUser) error {
	return r.db.WithContext(ctx).Create(matchMakerUser).Error
}

func (r *donutRepository) GetMatchMakerBySerial(ctx context.Context, serial string) (*MatchMakerEntity, error) {
	var matchMaker MatchMaker
	err := r.db.WithContext(ctx).Where("serial = ?", serial).First(&MatchMaker{}).Error
	if err != nil {
		return nil, err
	}
	return matchMaker.ToEntity(), nil
}
