package main

import (
	"context"
	"fmt"
	"math/rand"
)

type donutCall struct {
	People    People
	PeopleMap PeopleMap
	MatchMap  MatchMap

	repo DonutRepository
}

type DonutCall interface {
	Register(name string)
	Start()
	Pair()
	RePair()
	PairPeople(person1 *Person, person2 *Person)
	DoCall(person1 *Person, person2 *Person)
	AddPerson(name string)
	RemovePerson(name string)
	GetPerson(name string) *Person
	Print()
	CreateMatchMaker(ctx context.Context, options ...MatchMakerEntityOption) (string, error)
}

func NewDonutCall(donutRepository DonutRepository) DonutCall {
	return &donutCall{
		People:    make(People, 0),
		PeopleMap: make(map[string]bool),
		MatchMap:  make(map[string]string),
		repo:      donutRepository,
	}
}

func (dc *donutCall) Register(name string) {
	if dc.PeopleMap[name] {
		return
	}

	dc.People = append(dc.People, &Person{Name: name})
}

func (dc *donutCall) Start() {
	// now := time.Now()
	// if dc.StartAt.After(now) {
	// 	return
	// }

	dc.Pair()
}

func (dc *donutCall) Pair() {
	// round-robin balancing
	length := len(dc.People)

	people := make(People, length)
	copy(people, dc.People)

	for length >= 0 {
		if length == 0 {
			break
		}

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

		dc.PairPeople(person1, person2)
		length -= 2

		if p1Idx == p2Idx {
			continue
		}
	}
}

func (dc *donutCall) RePair() {
	// do the round-robin balancing again but only for people who have not been paired
	length := len(dc.GetRemaining())

	people := make(People, length)
	copy(people, dc.GetRemaining())

	for length >= 0 {
		if length == 0 {
			break
		}

		// special case, do a 3-way call
		if length == 1 {
			fmt.Println("3-way call", people[0].Name)
			break
		}

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

		dc.PairPeople(person1, person2)
		length -= 2

		if p1Idx == p2Idx {
			continue
		}
	}
}

func (dc *donutCall) PairPeople(person1 *Person, person2 *Person) {
	dc.PeopleMap[person1.Name] = false
	dc.PeopleMap[person2.Name] = false
	dc.MatchMap[person1.Name] = person2.Name
	dc.MatchMap[person2.Name] = person1.Name
	fmt.Println(person1.Name, "paired with", person2.Name)
}

func (dc *donutCall) DoCall(person1 *Person, person2 *Person) {
	fmt.Println(person1.Name, "called", person2.Name)
	dc.PeopleMap[person1.Name] = true
	dc.PeopleMap[person2.Name] = true
}

func (dc *donutCall) AddPerson(name string) {
	dc.People = append(dc.People, &Person{Name: name})
	fmt.Println("Added", name)
}

func (dc *donutCall) RemovePerson(name string) {
	for i, p := range dc.People {
		if p.Name == name {
			fmt.Println("Removed", name)
			delete(dc.PeopleMap, name)
			dc.People = append(dc.People[:i], dc.People[i+1:]...)
			break
		}
	}
}

func (dc *donutCall) GetPerson(name string) *Person {
	for _, p := range dc.People {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (dc *donutCall) Print() {
	fmt.Println("People:", dc.People.Print())
	fmt.Println("PeopleMap:", dc.PeopleMap)
	fmt.Println("MatchMap:", dc.MatchMap)
	fmt.Println("Completed:", dc.GetCompleted().Print())
	fmt.Println("Remaining:", dc.GetRemaining().Print())
}

func (dc *donutCall) GetCompleted() People {
	var completed People
	for _, p := range dc.People {
		if dc.PeopleMap[p.Name] {
			completed = append(completed, p)
		}
	}
	return completed
}

func (dc *donutCall) GetRemaining() People {
	var remaining People
	for _, p := range dc.People {
		if !dc.PeopleMap[p.Name] {
			remaining = append(remaining, p)
		}
	}
	return remaining
}

func (dc *donutCall) CreateMatchMaker(ctx context.Context, options ...MatchMakerEntityOption) (string, error) {
	matchMaker := &MatchMakerEntity{}
	matchMaker.ApplyOptions(options...)

	err := dc.repo.CreateMatchMaker(ctx, matchMaker)
	if err != nil {
		return "", err
	}

	return matchMaker.Serial, nil
}
