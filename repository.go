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
	UpdateMatchMakerStatusBySerial(ctx context.Context, serial string, status MatchMakerStatus) error

	UpdateSerialMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) error
	UpdateSerialMatchMakerUser(ctx context.Context, matchMakerUser *MatchMakerUserEntity) error
	UpdateStatusMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) error
	UpdateStatusMatchMakerUser(ctx context.Context, matchMakerUser *MatchMakerUserEntity) error
	DeleteMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) error

	GetMatchMakerBySerial(ctx context.Context, serial string) (*MatchMakerEntity, error)
	GetUsersByMatchMakerSerial(ctx context.Context, matchMakerSerial string) (MatchMakerUserEntities, error)
	GetUsersByMatchMakerSerialAndStatuses(ctx context.Context, matchMakerSerial string, status []MatchMakerUserStatus) (MatchMakerUserEntities, error)
	GetUsersByMatchMakerSerialAndUserReferences(ctx context.Context, matchMakerSerial string, userReferences []string) (MatchMakerUserEntities, error)
	GetUsersBySerial(ctx context.Context, serial string) (MatchMakerUserEntities, error)

	Database() *gorm.DB
}

func NewDonutRepository(db *gorm.DB) DonutRepository {
	return &donutRepository{
		db: db,
	}
}

func (r *donutRepository) Database() *gorm.DB {
	return r.db
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

// UpdateStatusMatchMakerUsers updates only for status of match maker users.
// This implementation is to handle for multiple update queries and transaction by gorm.
// The manager pattern is not necessary for this implementation.
func (r *donutRepository) UpdateStatusMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) (err error) {
	trx := r.db.WithContext(ctx).Begin()

	defer func() {
		if err != nil {
			err = trx.Rollback().Error
			return
		}

		err = trx.Commit().Error
		return
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
	q := fmt.Sprintf("%s = ? AND %s = ?", MatchMakerSerialColumn, UserReferenceColumn)
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, matchMakerUser := range matchMakerUsers {
			err := tx.WithContext(ctx).
				Where(q, matchMakerUser.MatchMakerSerial, matchMakerUser.UserReference).
				Delete(MatchMakerUsers{}).
				Error
			if err != nil {
				return err
			}
		}
		return nil
	})
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

func (r *donutRepository) GetUsersByMatchMakerSerialAndStatuses(ctx context.Context, matchMakerSerial string, status []MatchMakerUserStatus) (MatchMakerUserEntities, error) {
	var matchMakerUsers MatchMakerUsers
	q := fmt.Sprintf("%s = ? AND %s IN (?)", MatchMakerSerialColumn, StatusColumn)
	err := r.db.WithContext(ctx).Where(q, matchMakerSerial, status).Find(&matchMakerUsers).Error
	if err != nil {
		return nil, err
	}
	return matchMakerUsers.ToEntities(), nil
}

// UpdateSerialMatchMakerUsers updates only for serial and status of match maker users.
// This implementation is to handle for multiple update queries and transaction by gorm.
// The manager pattern is not necessary for this implementation.
func (r *donutRepository) UpdateSerialMatchMakerUsers(ctx context.Context, matchMakerUsers MatchMakerUserEntities) (err error) {
	trx := r.db.WithContext(ctx).Begin()

	defer func() {
		if err != nil {
			err = trx.Rollback().Error
			return
		}

		err = trx.Commit().Error
		return
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

func (r *donutRepository) UpdateSerialMatchMakerUser(ctx context.Context, matchMakerUser *MatchMakerUserEntity) error {
	q := fmt.Sprintf("%s = ? AND %s = ?", MatchMakerSerialColumn, UserReferenceColumn)
	updates := map[string]interface{}{
		SerialColumn: matchMakerUser.Serial,
		StatusColumn: MatchMakerUserStatusRunning,
	}
	return r.db.WithContext(ctx).
		Model(&MatchMakerUser{}).
		Where(q, matchMakerUser.MatchMakerSerial, matchMakerUser.UserReference).
		Updates(updates).
		Error
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

func (r *donutRepository) UpdateMatchMakerStatusBySerial(ctx context.Context, serial string, status MatchMakerStatus) error {
	q := fmt.Sprintf("%s = ?", SerialColumn)
	return r.db.WithContext(ctx).
		Model(&MatchMaker{}).
		Where(q, serial).
		Update(StatusColumn, status).
		Error
}

func (r *donutRepository) UpdateStatusMatchMakerUser(ctx context.Context, matchMakerUser *MatchMakerUserEntity) error {
	q := fmt.Sprintf("%s = ? AND %s = ?", MatchMakerSerialColumn, UserReferenceColumn)
	updates := map[string]interface{}{
		StatusColumn: matchMakerUser.Status,
	}
	return r.db.WithContext(ctx).
		Model(&MatchMakerUser{}).
		Where(q, matchMakerUser.MatchMakerSerial, matchMakerUser.UserReference).
		Updates(updates).
		Error
}
