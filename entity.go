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
	MatchMakerUserStatusRunning  MatchMakerUserStatus = "running"
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

func (m *MatchMakerEntity) Build(options ...MatchMakerEntityOption) *MatchMakerEntity {
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

	return m
}

func (m *MatchMakerEntity) Error() error {
	if m.Serial == "" {
		return fmt.Errorf("serial is empty")
	}

	if m.Name == "" {
		return fmt.Errorf("name is empty")
	}

	if m.StartTime.IsZero() {
		return fmt.Errorf("start time is zero")
	}

	if m.Duration == 0 {
		return fmt.Errorf("duration is zero")
	}

	return nil
}

type MatchMakerUserEntity struct {
	MatchMakerSerial string
	Serial           string
	UserReference    string
	Status           MatchMakerUserStatus
}

type MatchMakerUserEntityOption func(*MatchMakerUserEntity)

func WithMatchMakerUserEntityStatus(status MatchMakerUserStatus) MatchMakerUserEntityOption {
	return func(m *MatchMakerUserEntity) {
		m.Status = status
	}
}

func WithMatchMakerUserEntityUserReference(userReference string) MatchMakerUserEntityOption {
	return func(m *MatchMakerUserEntity) {
		m.UserReference = userReference
	}
}

func WithMatchMakerUserEntityMatchMakerSerial(matchMakerSerial string) MatchMakerUserEntityOption {
	return func(m *MatchMakerUserEntity) {
		m.MatchMakerSerial = matchMakerSerial
	}
}

func WithMatchMakerUserEntitySerial(serial string) MatchMakerUserEntityOption {
	return func(m *MatchMakerUserEntity) {
		m.Serial = serial
	}
}

func (m *MatchMakerUserEntity) Build(options ...MatchMakerUserEntityOption) *MatchMakerUserEntity {
	for _, opt := range options {
		opt(m)
	}

	if m.Status == "" {
		m.Status = MatchMakerUserStatusPending
	}

	return m
}

func (m *MatchMakerUserEntity) Error() error {
	if m.MatchMakerSerial == "" {
		return fmt.Errorf("match maker serial is empty")
	}

	if m.UserReference == "" {
		return fmt.Errorf("user reference is empty")
	}

	if m.Serial == "" {
		return fmt.Errorf("serial is empty")
	}

	return nil
}

type MatchMakerUserEntities []*MatchMakerUserEntity

func (m MatchMakerUserEntities) ToPeople() People {
	var people People
	for _, matchMakerUser := range m {
		if matchMakerUser == nil {
			continue
		}
		people = append(people, &Person{Name: matchMakerUser.UserReference})
	}
	return people
}

type MatchMakerInformation struct {
	MatchMaker *MatchMakerEntity
	Users      MatchMakerUserEntities
}
