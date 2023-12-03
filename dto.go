package main

import "time"

type MatchMakerStatus string

const (
	MatchMakerStatusPending  MatchMakerStatus = "pending"
	MatchMakerStatusRunning  MatchMakerStatus = "running"
	MatchMakerStatusFinished MatchMakerStatus = "finished"
)

type MatchMaker struct {
	ID          int64  `gorm:"primaryKey"`
	Serial      string `gorm:"uniqueIndex"`
	Name        string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

func (MatchMaker) TableName() string {
	return "matchmaker"
}

func (MatchMaker) FromEntity(entity *MatchMakerEntity) *MatchMaker {
	if entity == nil {
		return nil
	}

	return &MatchMaker{
		Serial:      entity.Serial,
		Name:        entity.Name,
		Description: entity.Description,
		StartTime:   entity.StartTime,
		EndTime:     entity.StartTime.Add(entity.Duration),
	}
}

func (m *MatchMaker) ToEntity() *MatchMakerEntity {
	if m == nil {
		return nil
	}

	return &MatchMakerEntity{
		Serial:      m.Serial,
		Name:        m.Name,
		Description: m.Description,
		StartTime:   m.StartTime,
		Duration:    m.EndTime.Sub(m.StartTime),
	}
}

type MatchMakerUser struct {
	ID               int64 `gorm:"primaryKey"`
	MatchMakerSerial string
	Serial           string `gorm:"uniqueIndex"`
	Username         string
	Status           MatchMakerStatus
}
