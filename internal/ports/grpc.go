package ports

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/edanko/users-api/internal/app"
	proto "github.com/edanko/users-api/pkg/genproto/v1"
)

type GrpcServer struct {
	app app.Application
}

func NewGrpcServer(application app.Application) GrpcServer {
	return GrpcServer{
		app: application,
	}
}

func (g GrpcServer) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (g GrpcServer) GetUsers(ctx context.Context, request *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (g GrpcServer) SyncUsers(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	// TODO implement me
	panic("implement me")
}

//
// func (g GrpcServer) GetKind(
// 	ctx context.Context,
// 	request *proto.GetKindRequest,
// ) (*proto.GetKindResponse, error) {
// 	id, err := uuid.Parse(request.GetId())
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
// 	appKind, err := g.app.Queries.GetKind.Handle(ctx, queries.GetKindRequest{
// 		ID: id,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	var status proto.KindStatus
// 	switch appKind.Status {
// 	case "published":
// 		status = proto.KindStatus_KIND_STATUS_PUBLISHED
// 	case "draft":
// 		status = proto.KindStatus_KIND_STATUS_DRAFT
// 	default:
// 		status = proto.KindStatus_KIND_STATUS_UNSPECIFIED
// 	}
//
// 	return &proto.GetKindResponse{
// 		Kind: &proto.Kind{
// 			Id:          appKind.ID.String(),
// 			Name:        appKind.Name,
// 			Description: appKind.Description,
// 			Status:      status,
// 		}}, nil
// }
//
// func (g GrpcServer) GetKindByName(
// 	ctx context.Context,
// 	request *proto.GetKindByNameRequest,
// ) (*proto.GetKindByNameResponse, error) {
// 	appKind, err := g.app.Queries.GetKindByName.Handle(ctx, queries.GetKindByNameRequest{
// 		Name: request.GetName(),
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	var status proto.KindStatus
// 	switch appKind.Status {
// 	case "published":
// 		status = proto.KindStatus_KIND_STATUS_PUBLISHED
// 	case "draft":
// 		status = proto.KindStatus_KIND_STATUS_DRAFT
// 	default:
// 		status = proto.KindStatus_KIND_STATUS_UNSPECIFIED
// 	}
//
// 	return &proto.GetKindByNameResponse{
// 		Kind: &proto.Kind{
// 			Id:          appKind.ID.String(),
// 			Name:        appKind.Name,
// 			Description: appKind.Description,
// 			Status:      status,
// 		}}, nil
// }
//
// func (g GrpcServer) ListKinds(
// 	ctx context.Context,
// 	request *proto.ListKindsRequest,
// ) (*proto.ListKindsResponse, error) {
// 	return nil, nil
// var limitInt int32
// var after string
// if request.Limit != nil {
// 	limitInt = *request.Limit
// }
// if request.After != nil {
// 	after = *request.After
// }
//
// // page := int32(math.Max(float64(pageInt), 1))
// limit := int32(math.Max(float64(limitInt), 20))
//
// // totalKinds, err := g.app.Queries.CountKinds.Handle(ctx, queries.CountKindsRequest{
// // 	Status: request.Status,
// // })
// // if err != nil {
// // 	return nil, status.Error(codes.Internal, err.Error())
// // }
// //
// // offset := (page * limit) - limit
//
// resp := &proto.ListKindsResponse{
// 	After: after,
// 	// Page:  page,
// 	Limit: limit,
// 	// Total: int32(totalKinds),
// }
//
// // if totalKinds < 1 || offset > (int32(totalKinds)-1) {
// // 	return resp, nil
// // }
//
// l := int(limit)
// // o := int(offset)
// a, err := uuid.Parse(after)
// if err != nil {
// 	return nil, status.Error(codes.Internal, err.Error())
// }
//
// var statusQuery *string
// if request.Status == nil {
// 	switch *request.Status {
// 	case proto.KindStatus_KIND_STATUS_DRAFT:
// 		*statusQuery = "draft"
// 	case proto.KindStatus_KIND_STATUS_PUBLISHED:
// 		*statusQuery = "published"
// 	}
// }
//
// appKinds, err := g.app.Queries.ListKinds.Handle(ctx, queries.ListKindsRequest{
// 	Status: statusQuery,
// 	Limit:  &l,
// 	After:  &a,
// 	// Offset: &o,
// })
// if err != nil {
// 	return nil, status.Error(codes.Internal, err.Error())
// }
//
// resp.Kinds = make([]*proto.Kind, 0, len(appKinds))
// for _, appKind := range appKinds {
// 	var appStatus proto.KindStatus
// 	switch appKind.Status {
// 	case "published":
// 		appStatus = proto.KindStatus_KIND_STATUS_PUBLISHED
// 	case "draft":
// 		appStatus = proto.KindStatus_KIND_STATUS_DRAFT
// 	default:
// 		appStatus = proto.KindStatus_KIND_STATUS_UNSPECIFIED
// 	}
//
// 	resp.Kinds = append(resp.Kinds, &proto.Kind{
// 		Id:          appKind.ID.String(),
// 		Login:        appKind.Login,
// 		Description: appKind.Description,
// 		Status:      appStatus,
// 	})
// }
// return resp, nil
// }

// func (g GrpcServer) CountKinds(
// 	ctx context.Context,
// 	request *proto.CountKindsRequest,
// ) (*proto.CountKindsResponse, error) {
// 	totalKinds, err := g.app.Queries.CountKinds.Handle(ctx, queries.CountKindsRequest{
// 		Status: request.Status,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	return &proto.CountKindsResponse{
// 		Total: int32(totalKinds),
// 	}, nil
// }
//
// func (g GrpcServer) ChangeKindDescription(
// 	ctx context.Context,
// 	request *proto.ChangeKindDescriptionRequest,
// ) (*proto.ChangeKindDescriptionResponse, error) {
// 	id, err := uuid.Parse(request.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	err = g.app.CommandBus.Send(ctx, commands.ChangeKindDescription{
// 		ID:          id,
// 		Description: request.Description,
// 	})
//
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	return &proto.ChangeKindDescriptionResponse{}, nil
// }
//
// func (g GrpcServer) ChangeKindName(
// 	ctx context.Context,
// 	request *proto.ChangeKindNameRequest,
// ) (*proto.ChangeKindNameResponse, error) {
// 	id, err := uuid.Parse(request.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	err = g.app.CommandBus.Send(ctx, commands.ChangeKindName{
// 		ID:   id,
// 		Name: request.Name,
// 	})
//
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	return &proto.ChangeKindNameResponse{}, nil
// }
//
// func (g GrpcServer) CreateKind(
// 	ctx context.Context,
// 	request *proto.CreateKindRequest,
// ) (*proto.CreateKindResponse, error) {
//
// 	var statusQuery string
//
// 	switch request.Status {
// 	case proto.KindStatus_KIND_STATUS_DRAFT:
// 		statusQuery = "draft"
// 	case proto.KindStatus_KIND_STATUS_PUBLISHED:
// 		statusQuery = "published"
// 	}
//
// 	err := g.app.CommandBus.Send(ctx, commands.CreateKind{
// 		ID:          uuid.New(),
// 		Name:        request.Name,
// 		Description: request.Description,
// 		Status:      statusQuery,
// 	})
//
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	return &proto.CreateKindResponse{}, nil
// }
//
// func (g GrpcServer) DeleteKind(
// 	ctx context.Context,
// 	request *proto.DeleteKindRequest,
// ) (*proto.DeleteKindResponse, error) {
// 	id, err := uuid.Parse(request.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	err = g.app.CommandBus.Send(ctx, commands.DeleteKind{
// 		ID: id,
// 	})
//
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	return &proto.DeleteKindResponse{}, nil
// }
//
// func (g GrpcServer) MakeKindDraft(
// 	ctx context.Context,
// 	request *proto.MakeKindDraftRequest,
// ) (*proto.MakeKindDraftResponse, error) {
// 	id, err := uuid.Parse(request.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	err = g.app.CommandBus.Send(ctx, commands.MakeKindDraft{
// 		ID: id,
// 	})
//
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	return &proto.MakeKindDraftResponse{}, nil
// }
//
// func (g GrpcServer) MakeKindPublished(
// 	ctx context.Context,
// 	request *proto.MakeKindPublishedRequest,
// ) (*proto.MakeKindPublishedResponse, error) {
// 	id, err := uuid.Parse(request.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	err = g.app.CommandBus.Send(ctx, commands.MakeKindPublished{
// 		ID: id,
// 	})
//
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}
//
// 	return &proto.MakeKindPublishedResponse{}, nil
// }
