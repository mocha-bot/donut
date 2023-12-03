package main

import (
	"fmt"
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

type PeopleMap map[string]bool

type MatchMap map[string]string

type MatchMakerEntity struct {
	Serial      string
	Name        string
	Description string
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

func (m *MatchMakerEntity) ApplyOptions(options ...MatchMakerEntityOption) {
	for _, opt := range options {
		opt(m)
	}

	if m.Serial == "" {
		m.Serial = fmt.Sprintf("%d", time.Now().UnixNano())
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
