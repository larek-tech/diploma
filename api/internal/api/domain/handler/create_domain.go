package handler

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	defaultTitlePattern = "%s (сценарий по умолчанию)"
	optimalTitlePattern = "%s (оптимальные параметры)"
)

// CreateDomain godoc
//
//	@Summary		Create new domain.
//	@Description	Creates new domain (group of sources used for RAG vector search).
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			req	body		pb.CreateDomainRequest	true	"Input data for creating domain"
//	@Success		201	{object}	pb.Domain				"Domain successfully created"
//	@Failure		400	{object}	string					"Failed to create domain"
//	@Router			/api/v1/domain/ [post]
func (h *Handler) CreateDomain(c *fiber.Ctx) error {
	var req pb.CreateDomainRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	resp, err := h.domainService.CreateDomain(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrCreateDomain, err.Error())
	}

	defaultScenarioReq := &pb.GetDefaultScenarioRequest{
		DefaultTitle: domainDefaultTitle(req.Title),
	}
	_, err = h.scenarioService.GetDefaultScenario(c.UserContext(), defaultScenarioReq)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return errs.WrapErr(shared.ErrGetScenario, err.Error())
		} else {
			params, err := h.mlService.GetDefaultParams(c.UserContext(), &emptypb.Empty{})
			if err != nil {
				return errs.WrapErr(shared.ErrGetScenario, err.Error())
			}

			createScenarioReq := &pb.CreateScenarioRequest{
				Title:        domainDefaultTitle(resp.GetTitle()),
				MultiQuery:   params.GetMultiQuery(),
				Reranker:     params.GetReranker(),
				VectorSearch: params.GetVectorSearch(),
				Model:        params.GetModel(),
				DomainId:     resp.GetId(),
			}
			scenario, err := h.scenarioService.CreateScenario(c.UserContext(), createScenarioReq)
			if err != nil {
				return errs.WrapErr(shared.ErrCreateScenario)
			}

			updateDomainReq := &pb.UpdateDomainRequest{
				DomainId:    resp.GetId(),
				Title:       resp.GetTitle(),
				SourceIds:   resp.GetSourceIds(),
				ScenarioIds: []int64{scenario.GetId()},
			}
			resp, err = h.domainService.UpdateDomain(c.UserContext(), updateDomainReq)
			if err != nil {
				return errs.WrapErr(shared.ErrUpdateDomain, err.Error())
			}

			go h.createOptimalScenario(c.UserContext(), resp)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func domainDefaultTitle(domainTitle string) string {
	return fmt.Sprintf(defaultTitlePattern, domainTitle)
}

func domainOptimalTitle(domainTitle string) string {
	return fmt.Sprintf(optimalTitlePattern, domainTitle)
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
