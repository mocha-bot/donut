package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Person struct {
	Name string
}

type People []*Person

func (p People) Get(i int) *Person {
	if i < 0 || i >= len(p) {
		return nil
	}
	return p[i]
}

func (p People) Print() string {
	var s string
	for i, person := range p {
		if i > 0 {
			s += ", "
		}
		s += person.Name
	}
	return s
}

type donutCall struct {
	People    People
	PeopleMap map[string]bool
	StartAt   time.Time
	Duration  time.Duration
	MatchMap  map[string]string
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
}

func NewDonutCall(startAt time.Time, duration time.Duration) DonutCall {
	return &donutCall{
		People:    make(People, 0),
		PeopleMap: make(map[string]bool),
		MatchMap:  make(map[string]string),
		StartAt:   startAt,
		Duration:  duration,
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

func main() {
	// create a slice of names
	names := []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Harry"}

	dc := NewDonutCall(time.Now(), 7*time.Minute)

	for _, name := range names {
		dc.Register(name)
	}

	fmt.Println()
	fmt.Println("Start...")

	dc.Start()

	fmt.Println()
	fmt.Println("Do calls...")

	dc.DoCall(dc.GetPerson("Alice"), dc.GetPerson("Bob"))
	dc.DoCall(dc.GetPerson("Charlie"), dc.GetPerson("David"))
	dc.DoCall(dc.GetPerson("Eve"), dc.GetPerson("Frank"))

	fmt.Println()
	fmt.Println("Add person...")

	dc.AddPerson("Ivan")
	dc.AddPerson("Goldi")
	dc.AddPerson("Samde")

	fmt.Println()
	fmt.Println("RePair...")

	dc.RePair()

	fmt.Println()
	fmt.Println("Remove person...")

	dc.RemovePerson("Ivan")

	fmt.Println()
	fmt.Println("RePair...")

	dc.RePair()

	fmt.Println()
	fmt.Println("Print...")

	dc.Print()
}
