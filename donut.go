package main

import (
	"context"
	"fmt"
	"math/rand"
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

	RegisterUsers(ctx context.Context, people MatchMakerUserEntities) error
	UnRegisterUsers(ctx context.Context, people MatchMakerUserEntities) error
}

func NewDonutCall(donutRepository DonutRepository) DonutCall {
	return &donutCall{
		repo: donutRepository,
	}
}

func (dc *donutCall) Call(ctx context.Context, matchMakerSerial string, people People) error {
	users, err := dc.repo.GetUsersByMatchMakerSerialAndUserReferences(ctx, matchMakerSerial, people.ToUserReferences())
	if err != nil {
		return err
	}

	matchMap := users.ToMatchMap()

	if len(matchMap) > 1 {
		return fmt.Errorf("wrong pair")
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

	if matchMaker.Status == MatchMakerStatusFinished {
		return fmt.Errorf("match maker is already finished")
	}

	return dc.Pair(ctx, matchMakerSerial)
}

func (dc *donutCall) Stop(ctx context.Context, matchMakerSerial string) error {
	return nil
}

func (dc *donutCall) GetPeople(ctx context.Context, matchMakerSerial string) (People, error) {
	matchMakerUsers, err := dc.repo.GetUsersByMatchMakerSerial(ctx, matchMakerSerial)
	return matchMakerUsers.ToPeople(), err
}

func (dc *donutCall) GetFinishedPeople(ctx context.Context, matchMakerSerial string) (People, error) {
	matchMakerUsers, err := dc.repo.GetUsersByMatchMakerSerialAndStatus(ctx, matchMakerSerial, MatchMakerUserStatusFinished)
	return matchMakerUsers.ToPeople(), err
}

func (dc *donutCall) GetPendingPeople(ctx context.Context, matchMakerSerial string) (People, error) {
	matchMakerUsers, err := dc.repo.GetUsersByMatchMakerSerialAndStatus(ctx, matchMakerSerial, MatchMakerUserStatusPending)
	return matchMakerUsers.ToPeople(), err
}

func (dc *donutCall) CreateMatchMaker(ctx context.Context, matchMaker *MatchMakerEntity) (string, error) {
	err := dc.repo.CreateMatchMaker(ctx, matchMaker)
	if err != nil {
		return "", err
	}

	return matchMaker.Serial, nil
}

func (dc *donutCall) RegisterUsers(ctx context.Context, people MatchMakerUserEntities) error {
	return dc.repo.CreateMatchMakerUsers(ctx, people)
}

func (dc *donutCall) UnRegisterUsers(ctx context.Context, people MatchMakerUserEntities) error {
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
	matchMakerUsers, err := dc.repo.GetUsersByMatchMakerSerialAndStatus(ctx, matchMakerSerial, MatchMakerUserStatusPending)
	if err != nil {
		return err
	}

	matchMakerPeople := matchMakerUsers.ToPeople()
	length := len(matchMakerPeople)

	people := make(People, length)
	copy(people, matchMakerPeople)

	matchMakerUsersEntities := make(MatchMakerUserEntities, 0)

	for length >= 0 {
		if length <= 0 {
			break
		}

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

	return dc.repo.UpdateSerialMatchMakerUsers(ctx, matchMakerUsersEntities)
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

	p1Idx := rand.Intn(length - 1)
	p2Idx := rand.Intn(length - 1)

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

	p1 := &MatchMakerUserEntity{}
	p1.Build(
		WithMatchMakerUserEntityMatchMakerSerial(matchMakerSerial),
		WithMatchMakerUserEntitySerial(matchmakingSerial),
		WithMatchMakerUserEntityUserReference(person1.Name),
		WithMatchMakerUserEntityStatus(MatchMakerUserStatusRunning),
	)

	p2 := &MatchMakerUserEntity{}
	p2.Build(
		WithMatchMakerUserEntityMatchMakerSerial(matchMakerSerial),
		WithMatchMakerUserEntitySerial(matchmakingSerial),
		WithMatchMakerUserEntityUserReference(person2.Name),
		WithMatchMakerUserEntityStatus(MatchMakerUserStatusRunning),
	)

	return append(matchMakerUsersEntities, p1, p2)
}
