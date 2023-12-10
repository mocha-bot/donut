package main

import (
	"context"

	donutv1 "buf.build/gen/go/mocha/remcall/protocolbuffers/go/donut/v1"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	svc DonutCall
}

func NewHandler(donutCall DonutCall) *Handler {
	return &Handler{
		svc: donutCall,
	}
}

func (h *Handler) CreateMatchMaker(ctx context.Context, req *connect.Request[donutv1.CreateMatchMakerRequest]) (*connect.Response[donutv1.CreateMatchMakerResponse], error) {
	serial, err := h.svc.CreateMatchMaker(ctx, parseCreateMatchMakerRequest(req))
	return parseCreateMatchMakerResponse(serial), err
}

func (h *Handler) GetMatchMakerInformation(ctx context.Context, req *connect.Request[donutv1.GetMatchMakerInformationRequest]) (*connect.Response[donutv1.GetMatchMakerInformationResponse], error) {
	info, err := h.svc.GetInformation(ctx, req.Msg.GetSerial())
	if err != nil {
		return nil, err
	}

	return parseGetMatchMakerInformationResponse(info), nil
}

func (h *Handler) StartMatchMaker(ctx context.Context, req *connect.Request[donutv1.StartMatchMakerRequest]) (*connect.Response[emptypb.Empty], error) {
	return connect.NewResponse(&emptypb.Empty{}), h.svc.Start(ctx, req.Msg.GetSerial())
}

func (h *Handler) StopMatchMaker(ctx context.Context, req *connect.Request[donutv1.StopMatchMakerRequest]) (*connect.Response[emptypb.Empty], error) {
	return connect.NewResponse(&emptypb.Empty{}), h.svc.Stop(ctx, req.Msg.GetSerial())
}

func (h *Handler) CallPeople(ctx context.Context, req *connect.Request[donutv1.CallPeopleRequest]) (*connect.Response[emptypb.Empty], error) {
	return connect.NewResponse(&emptypb.Empty{}), h.svc.Call(ctx, req.Msg.GetMatchmakerSerial(), parseCallPeopleRequest(req).ToPeople())
}

func (h *Handler) GetPeople(ctx context.Context, stream *connect.BidiStream[donutv1.GetPeopleRequest, donutv1.GetPeopleResponse]) error {
	for {
		msg, err := stream.Receive()
		if err != nil {
			return err
		}

		people, err := h.svc.GetPeople(ctx, msg.GetMatchmakerSerial())
		if err != nil {
			return err
		}

		err = stream.Send(parseGetPeopleResponse(people))
		if err != nil {
			return err
		}
	}
}

func (h *Handler) GetPeoplePair(ctx context.Context, req *connect.Request[donutv1.GetPeoplePairRequest]) (*connect.Response[donutv1.GetPeoplePairResponse], error) {
	matchMap, err := h.svc.GetPeoplePair(ctx, req.Msg.GetMatchmakerSerial())
	if err != nil {
		return nil, err
	}

	return parseGetPeoplePairResponse(matchMap), nil
}

func (h *Handler) RegisterPeople(ctx context.Context, stream *connect.BidiStream[donutv1.RegisterPeopleRequest, donutv1.RegisterPeopleResponse]) error {
	for {
		msg, err := stream.Receive()
		if err != nil {
			return err
		}

		err = h.svc.RegisterPeople(ctx, parseRegisterPeopleRequest(msg))
		if err != nil {
			return err
		}

		err = stream.Send(&donutv1.RegisterPeopleResponse{
			People: []*donutv1.Person{
				{
					Reference: msg.GetReference(),
				},
			},
		})
		if err != nil {
			return err
		}
	}
}

func (h *Handler) UnRegisterPeople(ctx context.Context, stream *connect.BidiStream[donutv1.UnRegisterPeopleRequest, donutv1.UnRegisterPeopleResponse]) error {
	for {
		msg, err := stream.Receive()
		if err != nil {
			return err
		}

		err = h.svc.UnRegisterPeople(ctx, parseUnRegisterPeopleRequest(msg))
		if err != nil {
			return err
		}

		err = stream.Send(&donutv1.UnRegisterPeopleResponse{
			People: []*donutv1.Person{
				{
					Reference: msg.GetReference(),
				},
			},
		})
		if err != nil {
			return err
		}
	}
}
