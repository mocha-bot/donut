package main

import (
	"time"

	donutv1 "buf.build/gen/go/mocha/remcall/protocolbuffers/go/donut/v1"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func parseCreateMatchMakerRequest(req *connect.Request[donutv1.CreateMatchMakerRequest]) *MatchMakerEntity {
	if req.Msg.MatchMaker == nil {
		return nil
	}

	matchMakerEntity := &MatchMakerEntity{}

	return matchMakerEntity.Build(
		WithMatchMakerEntityName(req.Msg.MatchMaker.GetName()),
		WithMatchMakerEntityDescription(req.Msg.MatchMaker.GetDescription()),
		WithMatchMakerEntityStartTime(req.Msg.MatchMaker.GetStartTime().AsTime()),
		WithMatchMakerEntityDuration(time.Duration(req.Msg.MatchMaker.GetDuration())),
	)
}

func parseCreateMatchMakerResponse(serial string) *connect.Response[donutv1.CreateMatchMakerResponse] {
	return connect.NewResponse(
		&donutv1.CreateMatchMakerResponse{
			Serial: serial,
		},
	)
}

func parseGetMatchMakerInformationResponse(info *MatchMakerInformation) *connect.Response[donutv1.GetMatchMakerInformationResponse] {
	return connect.NewResponse(
		&donutv1.GetMatchMakerInformationResponse{
			MatchMaker: &donutv1.MatchMaker{
				Serial:      info.MatchMaker.Serial,
				Name:        info.MatchMaker.Name,
				Description: info.MatchMaker.Description,
				StartTime:   timestamppb.New(info.MatchMaker.StartTime),
				Duration:    int32(info.MatchMaker.Duration),
			},
		},
	)
}

func parseRegisterPeopleRequest(req *donutv1.RegisterPeopleRequest) (parsed MatchMakerUserEntities) {
	entity := &MatchMakerUserEntity{}
	entity.Build(
		WithMatchMakerUserEntityMatchMakerSerial(req.GetMatchmakerSerial()),
		WithMatchMakerUserEntityUserReference(req.GetReference()),
		WithMatchMakerUserEntityStatus(MatchMakerUserStatusPending),
	)
	return append(parsed, entity)
}

func parseUnRegisterPeopleRequest(req *donutv1.UnRegisterPeopleRequest) (parsed MatchMakerUserEntities) {
	entity := &MatchMakerUserEntity{}
	return append(parsed, entity.Build(
		WithMatchMakerUserEntityMatchMakerSerial(req.GetMatchmakerSerial()),
		WithMatchMakerUserEntityUserReference(req.GetReference()),
	))
}

func parseCallPeopleRequest(req *connect.Request[donutv1.CallPeopleRequest]) (parsed MatchMakerUserEntities) {
	for _, person := range req.Msg.GetReferences() {
		entity := &MatchMakerUserEntity{}
		parsed = append(parsed, entity.Build(
			WithMatchMakerUserEntityMatchMakerSerial(req.Msg.GetMatchmakerSerial()),
			WithMatchMakerUserEntityUserReference(person),
		))
	}

	return parsed
}

func parseGetPeopleResponse(people People) *donutv1.GetPeopleResponse {
	peopleResp := make([]*donutv1.Person, len(people))

	for i, person := range people {
		peopleResp[i] = &donutv1.Person{
			Reference: person.Name,
		}
	}

	return &donutv1.GetPeopleResponse{
		People: peopleResp,
	}
}

func parseGetPeoplePairResponse(matchMap MatchMap) *connect.Response[donutv1.GetPeoplePairResponse] {
	peoplePairs := make([]*donutv1.PeoplePair, 0)

	for serial, people := range matchMap {
		peoplePairs = append(peoplePairs, &donutv1.PeoplePair{
			Serial: serial.String(),
			People: parseGetPeopleResponse(people).GetPeople(),
		})
	}

	return connect.NewResponse(
		&donutv1.GetPeoplePairResponse{
			PeoplePairs: peoplePairs,
		},
	)
}
