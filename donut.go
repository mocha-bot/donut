package main

import (
	"context"
)

type donutCall struct {
	repo DonutRepository
}

type DonutCall interface {
	// Pair()
	// RePair()

	CreateMatchMaker(ctx context.Context, matchMaker *MatchMakerEntity) (string, error)

	DoCall(ctx context.Context, matchMakerSerial string, people People) error
	Start(ctx context.Context, matchMakerSerial string) error

	GetInformation(ctx context.Context, matchMakerSerial string) (*MatchMakerInformation, error)

	GetPeople(ctx context.Context, matchMakerSerial string) (People, error)
	GetFinishedPeople(ctx context.Context, matchMakerSerial string) (People, error)
	GetPendingPeople(ctx context.Context, matchMakerSerial string) (People, error)

	RegisterUser(ctx context.Context, people MatchMakerUserEntities) error
	UnRegisterUser(ctx context.Context, people MatchMakerUserEntities) error
}

func NewDonutCall(donutRepository DonutRepository) DonutCall {
	return &donutCall{
		repo: donutRepository,
	}
}

// func (dc *donutCall) Pair() {
// 	// round-robin balancing
// 	length := len(dc.People)

// 	people := make(People, length)
// 	copy(people, dc.People)

// 	for length >= 0 {
// 		if length == 0 {
// 			break
// 		}

// 		p1Idx := rand.Intn(length - 1)
// 		p2Idx := rand.Intn(length - 1)

// 		// make sure p1 and p2 are not the same person
// 		for p1Idx == p2Idx {
// 			p2Idx = rand.Intn(length)
// 		}

// 		person1 := people[p1Idx]
// 		person2 := people[p2Idx]

// 		// remove paired people from the list of available people
// 		people = append(people[:p1Idx], people[p1Idx+1:]...)
// 		if p1Idx < p2Idx {
// 			p2Idx--
// 		}
// 		people = append(people[:p2Idx], people[p2Idx+1:]...)

// 		dc.PairPeople(person1, person2)
// 		length -= 2

// 		if p1Idx == p2Idx {
// 			continue
// 		}
// 	}
// }

// func (dc *donutCall) RePair() {
// 	// do the round-robin balancing again but only for people who have not been paired
// 	length := len(dc.GetRemaining())

// 	people := make(People, length)
// 	copy(people, dc.GetRemaining())

// 	for length >= 0 {
// 		if length == 0 {
// 			break
// 		}

// 		// special case, do a 3-way call
// 		if length == 1 {
// 			fmt.Println("3-way call", people[0].Name)
// 			break
// 		}

// 		p1Idx := rand.Intn(length - 1)
// 		p2Idx := rand.Intn(length - 1)

// 		// make sure p1 and p2 are not the same person
// 		for p1Idx == p2Idx {
// 			p2Idx = rand.Intn(length)
// 		}

// 		person1 := people[p1Idx]
// 		person2 := people[p2Idx]

// 		// remove paired people from the list of available people
// 		people = append(people[:p1Idx], people[p1Idx+1:]...)
// 		if p1Idx < p2Idx {
// 			p2Idx--
// 		}
// 		people = append(people[:p2Idx], people[p2Idx+1:]...)

// 		dc.PairPeople(person1, person2)
// 		length -= 2

// 		if p1Idx == p2Idx {
// 			continue
// 		}
// 	}
// }

func (dc *donutCall) DoCall(ctx context.Context, matchMakerSerial string, people People) error {
	// fmt.Println("Doing call", person1.Name, "with", person2.Name)

	return nil
}

func (dc *donutCall) Start(ctx context.Context, matchMakerSerial string) error {
	// dc.Pair()
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

func (dc *donutCall) RegisterUser(ctx context.Context, people MatchMakerUserEntities) error {
	return dc.repo.CreateMatchMakerUsers(ctx, people)
}

func (dc *donutCall) UnRegisterUser(ctx context.Context, people MatchMakerUserEntities) error {
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
