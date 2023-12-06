package main

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

type MatchMakerStatus string

const (
	MatchMakerStatusPending  MatchMakerStatus = "pending"
	MatchMakerStatusRunning  MatchMakerStatus = "running"
	MatchMakerStatusFinished MatchMakerStatus = "finished"
)

type MatchMakerUserStatus string

const (
	MatchMakerUserStatusPending  MatchMakerUserStatus = "pending"
	MatchMakerUserStatusFinished MatchMakerUserStatus = "finished"
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

type PeopleMap map[string]bool

type MatchMap map[string]string

type MatchMakerEntity struct {
	Serial      string
	Name        string
	Description string
	StartTime   time.Time
	Duration    time.Duration
}

func (m *MatchMakerEntity) GenerateSerial() {
	m.Serial = uuid.Must(uuid.NewV7()).String()
}

type MatchMakerEntityOption func(*MatchMakerEntity)

func WithMatchMakerEntityName(name string) MatchMakerEntityOption {
	return func(m *MatchMakerEntity) {
		m.Name = name
	}
}

func WithMatchMakerEntityDescription(description string) MatchMakerEntityOption {
	return func(m *MatchMakerEntity) {
		m.Description = description
	}
}

func WithMatchMakerEntityStartTime(startTime time.Time) MatchMakerEntityOption {
	return func(m *MatchMakerEntity) {
		m.StartTime = startTime
	}
}

func WithMatchMakerEntityDuration(duration time.Duration) MatchMakerEntityOption {
	return func(m *MatchMakerEntity) {
		m.Duration = duration
	}
}

func (m *MatchMakerEntity) Build(options ...MatchMakerEntityOption) {
	m.GenerateSerial()

	for _, opt := range options {
		opt(m)
	}

	if m.StartTime.IsZero() {
		m.StartTime = time.Now()
	}

	if m.Duration == 0 {
		m.Duration = 24 * time.Hour
	}

	if m.Name == "" {
		m.Name = fmt.Sprintf("MatchMaker-%s", m.Serial)
	}
}

type MatchMakerUserEntity struct {
	MatchMakerSerial string
	Serial           string
	Username         string
	Status           MatchMakerUserStatus
}

type MatchMakerUserEntities []*MatchMakerUserEntity

func (m MatchMakerUserEntities) ToPeople() People {
	var people People
	for _, matchMakerUser := range m {
		if matchMakerUser == nil {
			continue
		}
		people = append(people, &Person{Name: matchMakerUser.Username})
	}
	return people
}

type MatchMakerInformation struct {
	MatchMaker *MatchMakerEntity
	Users      MatchMakerUserEntities
}
