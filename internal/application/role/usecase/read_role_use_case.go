package role

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/role/mapper"
	shareddto "be-dashboard-nba/internal/application/shared/dto"
	"be-dashboard-nba/internal/domain/model"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
	"context"
	"math"

	"github.com/pkg/errors"
)

func (s *useCase) ReadRolesUseCase(
	ctx context.Context,
	args rolePresenter.ReadRolesRequest,
) (data model.RolePaginationResponse, err error) {
	r := s.newRoleRepo(s.db)

	params := mapper.ToReadRolesParams(&args)

	totalItems, err := r.ReadRolesCount(ctx, params)
	if err != nil {
		log(ctx).Error().Err(err).Msg("error query get role count")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(args.Limit)))
	hasNext := args.Page < totalPages
	hasPrev := args.Page > 1

	pagination := shareddto.Pagination{
		Page:       args.Page,
		PageSize:   args.Limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	rolesData, err := r.ReadRolesQuery(ctx, params)
	if err != nil {
		log(ctx).Error().Err(err).Msg("error query read roles")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data = model.RolePaginationResponse{
		Data:       rolesData,
		Pagination: pagination,
	}
	return
}
