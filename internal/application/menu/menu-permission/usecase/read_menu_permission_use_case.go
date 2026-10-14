package menu_permission

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/menu/menu-permission/mapper"
	shareddto "be-dashboard-nba/internal/application/shared/dto"
	"be-dashboard-nba/internal/domain/model"
	menuPermissionPresenter "be-dashboard-nba/internal/presentation/menu/menu-permission/presenter"
	"context"
	"database/sql"
	"math"

	"github.com/pkg/errors"
)

func (s *useCase) ReadMenuPermissionUseCase(
	ctx context.Context,
	args menuPermissionPresenter.ReadMenuPermissionListRequest,
	MenuID int,
) (data model.MenuPermissionPaginationResponse, err error) {

	mpr := s.newMenuPermissionRepo(s.db)
	mr := s.newMenuRepo(s.db)

	_, err = mr.ReadMenuByIDQuery(ctx, MenuID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Int("id", MenuID).Msg("menu detail not found")
			err = constant.ErrMenuIdNotFound
			return
		}
		log(ctx).Error().Err(err).Int("id", MenuID).Msg("error reading menu detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	params := mapper.ToReadMenuPermissionParams(&args, MenuID)

	totalItems, err := mpr.ReadMenuPermissionCount(ctx, params)
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

	menuPermissionData, err := mpr.ReadMenuPermissionListQuery(ctx, params)
	if err != nil {
		log(ctx).Error().Err(err).Msg("error query read list menu permission")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data = model.MenuPermissionPaginationResponse{
		Data:       menuPermissionData,
		Pagination: pagination,
	}

	return
}
