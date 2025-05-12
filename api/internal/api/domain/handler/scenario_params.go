package handler

import (
	"context"

	"github.com/larek-tech/diploma/api/internal/auth"
	authpb "github.com/larek-tech/diploma/api/internal/auth/pb"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) checkDefaultScenario(ctx context.Context, domain *pb.Domain, userID int64, roles []int64) (*pb.Domain, error) {
	defaultScenarioReq := &pb.GetDefaultScenarioRequest{
		DefaultTitle: domainDefaultTitle(domain.Title),
	}
	_, err := h.scenarioService.GetDefaultScenario(ctx, defaultScenarioReq)
	if err != nil {
		// internal error
		if status.Code(err) != codes.NotFound {
			return nil, errs.WrapErr(err)
		} else {
			// need to create default and optimal scenarios
			params, err := h.mlService.GetDefaultParams(ctx, &emptypb.Empty{})
			if err != nil {
				return nil, errs.WrapErr(err)
			}

			createScenarioReq := &pb.CreateScenarioRequest{
				Title:        domainDefaultTitle(domain.GetTitle()),
				MultiQuery:   params.GetMultiQuery(),
				Reranker:     params.GetReranker(),
				VectorSearch: params.GetVectorSearch(),
				Model:        params.GetModel(),
				DomainId:     domain.GetId(),
			}
			scenario, err := h.scenarioService.CreateScenario(ctx, createScenarioReq)
			if err != nil {
				return nil, errs.WrapErr(err)
			}

			updateDomainReq := &pb.UpdateDomainRequest{
				DomainId:    domain.GetId(),
				Title:       domain.GetTitle(),
				SourceIds:   domain.GetSourceIds(),
				ScenarioIds: []int64{scenario.GetId()},
			}
			resp, err := h.domainService.UpdateDomain(ctx, updateDomainReq)
			if err != nil {
				return nil, errs.WrapErr(err)
			}

			ctx = auth.PushUserMeta(ctx, &authpb.UserAuthMetadata{
				UserId: userID,
				Roles:  roles,
			})

			go h.createOptimalScenario(ctx, resp)
		}
	}
	return domain, nil
}

func (h *Handler) createOptimalScenario(ctx context.Context, domain *pb.Domain) {
	sourceIDsReq := &pb.GetSourceIDsRequest{
		SourceIds: domain.GetSourceIds(),
	}
	sourceIDs, err := h.sourceService.GetSourceIDs(ctx, sourceIDsReq)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get source ids")
		return
	}

	req := &pb.GetOptimalParamsRequest{
		SourceIds: sourceIDs.GetSourceIds(),
	}
	params, err := h.mlService.GetOptimalParams(ctx, req)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get optimal params")
		return
	}

	createScenarioReq := &pb.CreateScenarioRequest{
		Title:        domainOptimalTitle(domain.GetTitle()),
		MultiQuery:   params.GetMultiQuery(),
		Reranker:     params.GetReranker(),
		VectorSearch: params.GetVectorSearch(),
		Model:        params.GetModel(),
		DomainId:     domain.GetId(),
	}
	scenario, err := h.scenarioService.CreateScenario(ctx, createScenarioReq)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("create scenario")
		return
	}

	updateDomainReq := &pb.UpdateDomainRequest{
		DomainId:    domain.GetId(),
		Title:       domain.GetTitle(),
		SourceIds:   domain.GetSourceIds(),
		ScenarioIds: []int64{scenario.GetId()},
	}
	_, err = h.domainService.UpdateDomain(ctx, updateDomainReq)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("update domain")
		return
	}
}
