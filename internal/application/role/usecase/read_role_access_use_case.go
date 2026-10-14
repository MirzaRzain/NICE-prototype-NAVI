package role

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/role/mapper"
	shareddto "be-dashboard-nba/internal/application/shared/dto"
	"be-dashboard-nba/internal/domain/model"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
	"context"
	"database/sql"
	"math"

	"github.com/pkg/errors"
)

func (s *useCase) ReadRoleAccessUseCase(ctx context.Context, args rolePresenter.ReadRoleAccessesRequest, roleID int) (data model.RoleAccessPaginationResponse, err error) {
	r := s.newRoleRepo(s.db)
	params := mapper.ToReadRoleAccessParams(&args, roleID)

	totalItems, err := r.ReadRoleAccessCount(ctx, params)
	if err != nil {
		log(ctx).Error().Err(err).Msg("error query get menu permission count")
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

	_, err = r.ReadRoleByIDQuery(ctx, roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Int("Id", roleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		log(ctx).Error().Err(err).Int("id", roleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	roleAccessData, err := r.ReadRoleAccessQuery(ctx, params)
	if err != nil {
		log(ctx).Error().Err(err).Int("id", roleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data = model.RoleAccessPaginationResponse{
		Data:       roleAccessData,
		Pagination: pagination,
	}
	return
}
