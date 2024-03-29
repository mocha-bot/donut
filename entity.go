package main

import (
	"fmt"
	"time"
)

type MatchMakerStatus string

const (
	MatchMakerStatusPending  MatchMakerStatus = "pending"
	MatchMakerStatusRunning  MatchMakerStatus = "running"
	MatchMakerStatusFinished MatchMakerStatus = "finished"
	MatchMakerStatusStopped  MatchMakerStatus = "stopped"
)

type MatchMakerUserStatus string

const (
	MatchMakerUserStatusPending  MatchMakerUserStatus = "pending"
	MatchMakerUserStatusRunning  MatchMakerUserStatus = "running"
	MatchMakerUserStatusFinished MatchMakerUserStatus = "finished"
	MatchMakerUserStatusStopped  MatchMakerUserStatus = "stopped"
)

type Person struct {
	Name string
}

type People []*Person

func (p People) ToUserReferences() []string {
	var userReferences []string
	for _, person := range p {
		if person == nil {
			continue
		}
		userReferences = append(userReferences, person.Name)
	}
	return userReferences
}

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

type MatchMakerUserSerial string

func (m MatchMakerUserSerial) String() string {
	return string(m)
}

type MatchMap map[MatchMakerUserSerial]People

func (m MatchMap) First() (MatchMakerUserSerial, People) {
	for serial, match := range m {
		return serial, match
	}
	return "", nil
}

type MatchMakerEntity struct {
	Serial      string
	Name        string
	Description string
	Status      MatchMakerStatus
	StartTime   time.Time
	Duration    time.Duration
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
	m.Serial = GenerateSerial()
	m.Status = MatchMakerStatusPending

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

func (m MatchMakerUserEntities) ToMatchMap() MatchMap {
	matchMap := make(MatchMap)
	for _, matchMakerUser := range m {
		if matchMakerUser == nil {
			continue
		}
		match, ok := matchMap[MatchMakerUserSerial(matchMakerUser.Serial)]
		if !ok {
			matchMap[MatchMakerUserSerial(matchMakerUser.Serial)] = make(People, 0)
		}
		match = append(match, &Person{Name: matchMakerUser.UserReference})
		matchMap[MatchMakerUserSerial(matchMakerUser.Serial)] = match
	}
	return matchMap
}

type MatchMakerInformation struct {
	MatchMaker *MatchMakerEntity
	Users      MatchMakerUserEntities
	Pairs      MatchMap
}
