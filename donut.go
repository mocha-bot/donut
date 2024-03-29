package main

import (
	"context"
	"fmt"
	"math/rand"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
)

type donutCall struct {
	repo DonutRepository
}

type DonutCall interface {
	Start(ctx context.Context, matchMakerSerial string) error
	Stop(ctx context.Context, matchMakerSerial string) error

	Pair(ctx context.Context, matchMakerSerial string) error
	Call(ctx context.Context, matchMakerSerial string, people People) error

	CreateMatchMaker(ctx context.Context, matchMaker *MatchMakerEntity) (string, error)

	GetInformation(ctx context.Context, matchMakerSerial string) (*MatchMakerInformation, error)

	GetPeople(ctx context.Context, matchMakerSerial string) (People, error)
	GetFinishedPeople(ctx context.Context, matchMakerSerial string) (People, error)
	GetPendingPeople(ctx context.Context, matchMakerSerial string) (People, error)
	GetPeoplePair(ctx context.Context, matchMakerSerial string) (MatchMap, error)

	RegisterPeople(ctx context.Context, people MatchMakerUserEntities) error
	UnRegisterPeople(ctx context.Context, people MatchMakerUserEntities) error
}

func NewDonutCall(donutRepository DonutRepository) DonutCall {
	return &donutCall{
		repo: donutRepository,
	}
}

func (dc *donutCall) Call(ctx context.Context, matchMakerSerial string, people People) error {
	matchMaker, err := dc.repo.GetMatchMakerBySerial(ctx, matchMakerSerial)
	if err != nil {
		return err
	}

	if matchMaker.Status != MatchMakerStatusRunning {
		return fmt.Errorf("match maker is not running")
	}

	users, err := dc.repo.GetUsersByMatchMakerSerialAndUserReferences(ctx, matchMakerSerial, people.ToUserReferences())
	if err != nil {
		return err
	}

	matchMap := users.ToMatchMap()

	if len(matchMap) > 1 {
		return fmt.Errorf("expected one pair but found %d", len(matchMap))
	}

	matchMakerUserSerial, _ := matchMap.First()
	usersRegistered, err := dc.repo.GetUsersBySerial(ctx, matchMakerUserSerial.String())
	if err != nil {
		return err
	}

	if len(usersRegistered) != len(people) {
		return fmt.Errorf("pair is lack of people")
	}

	matchMakerUsersEntities := make(MatchMakerUserEntities, 0)

	for _, person := range people {
		if person == nil {
			continue
		}

		matchMakerUser := &MatchMakerUserEntity{}
		matchMakerUser.Build(
			WithMatchMakerUserEntityMatchMakerSerial(matchMakerSerial),
			WithMatchMakerUserEntitySerial(matchMakerUserSerial.String()),
			WithMatchMakerUserEntityUserReference(person.Name),
			WithMatchMakerUserEntityStatus(MatchMakerUserStatusFinished),
		)

		matchMakerUsersEntities = append(matchMakerUsersEntities, matchMakerUser)
	}

	return dc.repo.UpdateStatusMatchMakerUsers(ctx, matchMakerUsersEntities)
}

func (dc *donutCall) Start(ctx context.Context, matchMakerSerial string) error {
	matchMaker, err := dc.repo.GetMatchMakerBySerial(ctx, matchMakerSerial)
	if err != nil {
		return err
	}

	if matchMaker.Status == MatchMakerStatusRunning {
		return fmt.Errorf("match maker is already running")
	}

	if matchMaker.Status == MatchMakerStatusFinished || matchMaker.Status == MatchMakerStatusStopped {
		return fmt.Errorf("match maker is already finished or stopped")
	}

	return dc.Pair(ctx, matchMakerSerial)
}

func (dc *donutCall) Stop(ctx context.Context, matchMakerSerial string) error {
	matchMaker, err := dc.repo.GetMatchMakerBySerial(ctx, matchMakerSerial)
	if err != nil {
		return err
	}

	if matchMaker.Status == MatchMakerStatusFinished {
		return nil
	}

	matchMakerUsers, err := dc.repo.GetUsersByMatchMakerSerialAndStatuses(ctx, matchMakerSerial, []MatchMakerUserStatus{MatchMakerUserStatusPending, MatchMakerUserStatusRunning})
	if err != nil {
		return err
	}

	trManagerSettingOptions, err := settings.New(settings.WithPropagation(trm.PropagationRequired))
	if err != nil {
		return err
	}

	trManagerSetting, err := trmgorm.NewSettings(trManagerSettingOptions)
	if err != nil {
		return err
	}

	trManager, err := manager.New(trmgorm.NewDefaultFactory(dc.repo.Database()), manager.WithSettings(trManagerSetting))
	if err != nil {
		return err
	}

	return trManager.Do(ctx, func(ctx context.Context) error {
		for _, matchMakerUser := range matchMakerUsers {
			if matchMakerUser == nil {
				continue
			}
			matchMakerUser.Status = MatchMakerUserStatusStopped
			err := dc.repo.UpdateStatusMatchMakerUser(ctx, matchMakerUser)
			if err != nil {
				return err
			}
		}
		return dc.repo.UpdateMatchMakerStatusBySerial(ctx, matchMakerSerial, MatchMakerStatusFinished)
	})
}

func (dc *donutCall) GetPeople(ctx context.Context, matchMakerSerial string) (People, error) {
	matchMakerUsers, err := dc.repo.GetUsersByMatchMakerSerial(ctx, matchMakerSerial)
	return matchMakerUsers.ToPeople(), err
}

func (dc *donutCall) GetFinishedPeople(ctx context.Context, matchMakerSerial string) (People, error) {
	matchMakerUsers, err := dc.repo.GetUsersByMatchMakerSerialAndStatuses(ctx, matchMakerSerial, []MatchMakerUserStatus{MatchMakerUserStatusFinished})
	return matchMakerUsers.ToPeople(), err
}

func (dc *donutCall) GetPendingPeople(ctx context.Context, matchMakerSerial string) (People, error) {
	matchMakerUsers, err := dc.repo.GetUsersByMatchMakerSerialAndStatuses(ctx, matchMakerSerial, []MatchMakerUserStatus{MatchMakerUserStatusPending})
	return matchMakerUsers.ToPeople(), err
}

func (dc *donutCall) CreateMatchMaker(ctx context.Context, matchMaker *MatchMakerEntity) (string, error) {
	err := dc.repo.CreateMatchMaker(ctx, matchMaker)
	if err != nil {
		return "", err
	}

	return matchMaker.Serial, nil
}

func (dc *donutCall) RegisterPeople(ctx context.Context, people MatchMakerUserEntities) error {
	return dc.repo.CreateMatchMakerUsers(ctx, people)
}

func (dc *donutCall) UnRegisterPeople(ctx context.Context, people MatchMakerUserEntities) error {
	return dc.repo.DeleteMatchMakerUsers(ctx, people)
}

func (dc *donutCall) GetInformation(ctx context.Context, matchMakerSerial string) (*MatchMakerInformation, error) {
	matchMaker, err := dc.repo.GetMatchMakerBySerial(ctx, matchMakerSerial)
	if err != nil {
		return nil, err
	}

	matchMakerUsers, err := dc.repo.GetUsersByMatchMakerSerial(ctx, matchMakerSerial)
	if err != nil {
		return nil, err
	}

	return &MatchMakerInformation{
		MatchMaker: matchMaker,
		Users:      matchMakerUsers,
	}, nil
}

func (dc *donutCall) Pair(ctx context.Context, matchMakerSerial string) error {
	matchMakerUsers, err := dc.repo.GetUsersByMatchMakerSerialAndStatuses(ctx, matchMakerSerial, []MatchMakerUserStatus{MatchMakerUserStatusPending})
	if err != nil {
		return err
	}

	matchMakerPeople := matchMakerUsers.ToPeople()
	length := len(matchMakerPeople)

	people := make(People, length)
	copy(people, matchMakerPeople)

	matchMakerUsersEntities := make(MatchMakerUserEntities, 0)

	for length > 0 {
		switch {
		case length == 3:
			people, matchMakerUsersEntities = processThreeWayCall(people, matchMakerSerial, matchMakerUsersEntities)
			length = 0
		case length == 1:
			people = processOneWayCall(people)
			length = 0
		default:
			matchMakerUsersEntities = processTwoWayCall(length, people, matchMakerSerial, matchMakerUsersEntities)
			length -= 2
		}
	}

	trManagerSettingOptions, err := settings.New(settings.WithPropagation(trm.PropagationRequired))
	if err != nil {
		return err
	}

	trManagerSetting, err := trmgorm.NewSettings(trManagerSettingOptions)
	if err != nil {
		return err
	}

	trManager, err := manager.New(trmgorm.NewDefaultFactory(dc.repo.Database()), manager.WithSettings(trManagerSetting))
	if err != nil {
		return err
	}

	return trManager.Do(ctx, func(ctx context.Context) error {
		for _, matchMakerUser := range matchMakerUsersEntities {
			if matchMakerUser == nil {
				continue
			}
			err := dc.repo.UpdateSerialMatchMakerUser(ctx, matchMakerUser)
			if err != nil {
				return err
			}
		}
		return dc.repo.UpdateMatchMakerStatusBySerial(ctx, matchMakerSerial, MatchMakerStatusRunning)
	})
}

func (dc *donutCall) GetPeoplePair(ctx context.Context, matchMakerSerial string) (MatchMap, error) {
	matchMakerUsers, err := dc.repo.GetUsersByMatchMakerSerial(ctx, matchMakerSerial)
	if err != nil {
		return nil, err
	}

	return matchMakerUsers.ToMatchMap(), nil
}

func processOneWayCall(people People) People {
	fmt.Printf("There is one person left: %s\n", people[0].Name)
	people = people[:0]
	return people
}

func processThreeWayCall(people People, matchMakerSerial string, matchMakerUsersEntities MatchMakerUserEntities) (People, MatchMakerUserEntities) {
	matchmakingSerial := GenerateSerial()

	for _, person := range people {
		if person == nil {
			continue
		}

		matchMakerUser := &MatchMakerUserEntity{}
		matchMakerUser.Build(
			WithMatchMakerUserEntityMatchMakerSerial(matchMakerSerial),
			WithMatchMakerUserEntitySerial(matchmakingSerial),
			WithMatchMakerUserEntityUserReference(person.Name),
			WithMatchMakerUserEntityStatus(MatchMakerUserStatusRunning),
		)

		matchMakerUsersEntities = append(matchMakerUsersEntities, matchMakerUser)
	}
	people = people[:0]
	return people, matchMakerUsersEntities
}

func processTwoWayCall(length int, people People, matchMakerSerial string, matchMakerUsersEntities MatchMakerUserEntities) MatchMakerUserEntities {
	if length == 0 {
		return matchMakerUsersEntities
	}

	n := length - 1

	p1Idx := rand.Intn(n)
	p2Idx := rand.Intn(n)

	for p1Idx == p2Idx {
		p2Idx = rand.Intn(length)
	}

	person1 := people[p1Idx]
	person2 := people[p2Idx]

	people = append(people[:p1Idx], people[p1Idx+1:]...)
	if p1Idx < p2Idx {
		p2Idx--
	}
	people = append(people[:p2Idx], people[p2Idx+1:]...)

	matchmakingSerial := GenerateSerial()

	for _, person := range []*Person{person1, person2} {
		if person == nil {
			continue
		}

		matchMakerUser := &MatchMakerUserEntity{}
		matchMakerUser.Build(
			WithMatchMakerUserEntityMatchMakerSerial(matchMakerSerial),
			WithMatchMakerUserEntitySerial(matchmakingSerial),
			WithMatchMakerUserEntityUserReference(person.Name),
			WithMatchMakerUserEntityStatus(MatchMakerUserStatusRunning),
		)

		matchMakerUsersEntities = append(matchMakerUsersEntities, matchMakerUser)
	}

	return matchMakerUsersEntities
}
