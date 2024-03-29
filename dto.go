package main

import "time"

const (
	MatchMakerSerialColumn = "matchmaker_serial"
	UserReferenceColumn    = "user_reference"
	SerialColumn           = "serial"
	StatusColumn           = "status"
)

type MatchMaker struct {
	Serial      string `gorm:"uniqueIndex"`
	Name        string
	Description string
	Status      MatchMakerStatus
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
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
		Status:      entity.Status,
		StartTime:   entity.StartTime,
		EndTime:     entity.StartTime.Add(entity.Duration * Day),
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
		Status:      m.Status,
	}
}

type MatchMakerUser struct {
	MatchMakerSerial string `gorm:"column:matchmaker_serial"`
	Serial           string `gorm:"uniqueIndex"`
	UserReference    string
	Status           MatchMakerUserStatus
	DeletedAt        *time.Time
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

func (MatchMakerUser) TableName() string {
	return "matchmaker_user"
}

func (m *MatchMakerUser) FromEntity(entity *MatchMakerUserEntity) *MatchMakerUser {
	if entity == nil {
		return nil
	}

	return &MatchMakerUser{
		MatchMakerSerial: entity.MatchMakerSerial,
		Serial:           entity.Serial,
		UserReference:    entity.UserReference,
		Status:           entity.Status,
	}
}

func (m *MatchMakerUser) ToEntity() *MatchMakerUserEntity {
	if m == nil {
		return nil
	}

	return &MatchMakerUserEntity{
		MatchMakerSerial: m.MatchMakerSerial,
		Serial:           m.Serial,
		UserReference:    m.UserReference,
		Status:           m.Status,
	}
}

type MatchMakerUsers []*MatchMakerUser

func (m MatchMakerUsers) FromEntities(entities MatchMakerUserEntities) MatchMakerUsers {
	var matchMakerUsers MatchMakerUsers
	for _, entity := range entities {
		if entity == nil {
			continue
		}
		matchMakerUsers = append(matchMakerUsers, (&MatchMakerUser{}).FromEntity(entity))
	}
	return matchMakerUsers
}

func (m MatchMakerUsers) ToEntities() MatchMakerUserEntities {
	var entities MatchMakerUserEntities
	for _, matchMakerUser := range m {
		if matchMakerUser == nil {
			continue
		}
		entities = append(entities, matchMakerUser.ToEntity())
	}
	return entities
}
