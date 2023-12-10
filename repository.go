package main

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type donutRepository struct {
	db *gorm.DB
}

type DonutRepository interface {
	CreateMatchMaker(ctx context.Context, matchMaker *MatchMakerEntity) error
	CreateMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) error

	UpdateSerialMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) error
	UpdateStatusMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) error
	DeleteMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) error

	GetMatchMakerBySerial(ctx context.Context, serial string) (*MatchMakerEntity, error)
	GetUsersByMatchMakerSerial(ctx context.Context, matchMakerSerial string) (MatchMakerUserEntities, error)
	GetUsersByMatchMakerSerialAndStatus(ctx context.Context, matchMakerSerial string, status MatchMakerUserStatus) (MatchMakerUserEntities, error)
	GetUsersByMatchMakerSerialAndUserReferences(ctx context.Context, matchMakerSerial string, userReferences []string) (MatchMakerUserEntities, error)
	GetUsersBySerial(ctx context.Context, serial string) (MatchMakerUserEntities, error)
}

func NewDonutRepository(db *gorm.DB) DonutRepository {
	return &donutRepository{
		db: db,
	}
}

func (r *donutRepository) CreateMatchMaker(ctx context.Context, matchMaker *MatchMakerEntity) error {
	return r.db.WithContext(ctx).Create(MatchMaker{}.FromEntity(matchMaker)).Error
}

func (r *donutRepository) CreateMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) error {
	clauses := clause.OnConflict{DoNothing: true}
	return r.db.
		WithContext(ctx).
		Clauses(clauses).
		Model(&MatchMakerUser{}).
		Create(MatchMakerUsers{}.FromEntities(matchMakerUsers)).
		Error
}

func (r *donutRepository) UpdateStatusMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) (err error) {
	trx := r.db.WithContext(ctx).Begin()

	defer func() {
		if err != nil {
			trx.Rollback()
			return
		}

		err := trx.Commit()
		if err != nil {
			return
		}
	}()

	for _, matchMakerUser := range matchMakerUsers {
		if matchMakerUser == nil {
			continue
		}
		q := fmt.Sprintf("%s = ?", SerialColumn)
		err := trx.Model(&MatchMakerUser{}).
			Where(q, matchMakerUser.Serial).
			Update(StatusColumn, matchMakerUser.Status).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *donutRepository) DeleteMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) error {
	return r.db.WithContext(ctx).Delete(MatchMakerUsers{}.FromEntities(matchMakerUsers)).Error
}

func (r *donutRepository) GetMatchMakerBySerial(ctx context.Context, serial string) (*MatchMakerEntity, error) {
	var matchMaker MatchMaker
	q := fmt.Sprintf("%s = ?", SerialColumn)
	err := r.db.WithContext(ctx).Where(q, serial).First(&matchMaker).Error
	if err != nil {
		return nil, err
	}
	return matchMaker.ToEntity(), nil
}

func (r *donutRepository) GetUsersByMatchMakerSerial(ctx context.Context, matchMakerSerial string) (MatchMakerUserEntities, error) {
	var matchMakerUsers MatchMakerUsers
	q := fmt.Sprintf("%s = ?", MatchMakerSerialColumn)
	err := r.db.WithContext(ctx).Where(q, matchMakerSerial).Find(&matchMakerUsers).Error
	if err != nil {
		return nil, err
	}
	return matchMakerUsers.ToEntities(), nil
}

func (r *donutRepository) GetUsersByMatchMakerSerialAndStatus(ctx context.Context, matchMakerSerial string, status MatchMakerUserStatus) (MatchMakerUserEntities, error) {
	var matchMakerUsers MatchMakerUsers
	q := fmt.Sprintf("%s = ? AND %s = ?", MatchMakerSerialColumn, StatusColumn)
	err := r.db.WithContext(ctx).Where(q, matchMakerSerial, status).Find(&matchMakerUsers).Error
	if err != nil {
		return nil, err
	}
	return matchMakerUsers.ToEntities(), nil
}

func (r *donutRepository) UpdateSerialMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) (err error) {
	trx := r.db.WithContext(ctx).Begin()

	defer func() {
		if err != nil {
			trx.Rollback()
			return
		}

		trx.Commit()
	}()

	for _, matchMakerUser := range matchMakerUsers {
		if matchMakerUser == nil {
			continue
		}
		q := fmt.Sprintf("%s = ? AND %s = ?", MatchMakerSerialColumn, UserReferenceColumn)
		updates := map[string]interface{}{
			SerialColumn: matchMakerUser.Serial,
			StatusColumn: MatchMakerUserStatusRunning,
		}

		err := trx.Model(&MatchMakerUser{}).
			Where(q, matchMakerUser.MatchMakerSerial, matchMakerUser.UserReference).
			Updates(updates).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *donutRepository) GetUsersByMatchMakerSerialAndUserReferences(ctx context.Context, matchMakerSerial string, userReferences []string) (MatchMakerUserEntities, error) {
	const batchSize = 1000
	var allMatchMakerUsers MatchMakerUserEntities

	for i := 0; i < len(userReferences); i += batchSize {
		end := i + batchSize
		if end > len(userReferences) {
			end = len(userReferences)
		}
		var batchMatchMakerUsers MatchMakerUsers
		q := fmt.Sprintf("%s = ? AND %s IN ?", MatchMakerSerialColumn, UserReferenceColumn)
		err := r.db.WithContext(ctx).Where(q, matchMakerSerial, userReferences[i:end]).Find(&batchMatchMakerUsers).Error
		if err != nil {
			return nil, err
		}
		allMatchMakerUsers = append(allMatchMakerUsers, batchMatchMakerUsers.ToEntities()...)
	}

	return allMatchMakerUsers, nil
}

func (r *donutRepository) GetUsersBySerial(ctx context.Context, serial string) (MatchMakerUserEntities, error) {
	var matchMakerUsers MatchMakerUsers
	q := fmt.Sprintf("%s = ?", SerialColumn)
	err := r.db.WithContext(ctx).Where(q, serial).Find(&matchMakerUsers).Error
	if err != nil {
		return nil, err
	}
	return matchMakerUsers.ToEntities(), nil
}
