package main

import (
	"context"
	"math/rand"
)

type donutCall struct {
	repo DonutRepository
}

type DonutCall interface {
	CreateMatchMaker(ctx context.Context, matchMaker *MatchMakerEntity) (string, error)

	Pair(ctx context.Context, matchMakerSerial string) error
	RePair(ctx context.Context, matchMakerSerial string) error

	DoCall(ctx context.Context, matchMakerSerial string, people People) error
	Start(ctx context.Context, matchMakerSerial string) error
	Stop(ctx context.Context, matchMakerSerial string) error

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

func (dc *donutCall) DoCall(ctx context.Context, matchMakerSerial string, people People) error {
	// fmt.Println("Doing call", person1.Name, "with", person2.Name)

	return nil
}

func (dc *donutCall) Start(ctx context.Context, matchMakerSerial string) error {
	return dc.Pair(ctx, matchMakerSerial)
}

func (dc *donutCall) Stop(ctx context.Context, matchMakerSerial string) error {
	return dc.RePair(ctx, matchMakerSerial)
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
		if length == 0 {
			break
		}

		// get random person
		p1Idx := rand.Intn(length - 1)
		p2Idx := rand.Intn(length - 1)

		// make sure p1 and p2 are not the same person
		for p1Idx == p2Idx {
			p2Idx = rand.Intn(length)
		}

		person1 := people[p1Idx]
		person2 := people[p2Idx]

		// remove paired people from the list of available people
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

		matchMakerUsersEntities = append(matchMakerUsersEntities, p1, p2)

		length -= 2
		if p1Idx == p2Idx {
			continue
		}
	}

	return dc.repo.UpdateSerialMatchMakerUsers(ctx, matchMakerUsersEntities)
}

func (dc *donutCall) RePair(ctx context.Context, matchMakerSerial string) error {
	return nil
}
