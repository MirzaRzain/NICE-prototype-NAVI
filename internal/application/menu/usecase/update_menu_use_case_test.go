package menu

import (
	"context"
	"database/sql"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/menu/dto"
	"be-dashboard-nba/internal/domain/model"
	menuPresenter "be-dashboard-nba/internal/presentation/menu/presenter"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestUpdateMenuUseCase(t *testing.T) {
	parentID := 100
	active := true
	display := true

	tests := []struct {
		name    string
		menuID  int
		req     menuPresenter.UpdateMenuRequest
		setup   func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock)
		wantErr error
	}{
		{
			name:   "negative/not-found",
			menuID: 999,
			req: menuPresenter.UpdateMenuRequest{
				Name:    "X",
				Group:   "main",
				Active:  &active,
				Display: &display,
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()
				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 999).Return(model.Menu{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrMenuIdNotFound,
		},
		{
			name:   "negative/count-children-error",
			menuID: 1,
			req: menuPresenter.UpdateMenuRequest{
				Name:    "X",
				Group:   "main",
				Active:  &active,
				Display: &display,
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()
				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 1).Return(model.Menu{ID: 1}, nil).Once()
				menuRepo.On("CountMenuChildren", testifymock.Anything, 1).Return(0, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name:   "edge/has-children-parent-changed-forbidden",
			menuID: 1,
			req: menuPresenter.UpdateMenuRequest{
				ParentID: &parentID,
				Name:     "X",
				Group:    "main",
				Active:   &active,
				Display:  &display,
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 1).Return(model.Menu{ID: 1}, nil).Once()
				menuRepo.On("CountMenuChildren", testifymock.Anything, 1).Return(1, nil).Once()
			},
			wantErr: constant.ErrMenuHasChildren,
		},
		{
			name:   "positive/move-to-parent-recomputes-sort-and-inherits-group",
			menuID: 1,
			req: menuPresenter.UpdateMenuRequest{
				ParentID: &parentID,
				Name:     "X",
				Group:    "main",
				Active:   &active,
				Display:  &display,
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 1).Return(model.Menu{
					ID:       1,
					Group:    "main",
					ParentID: sql.NullInt32{Valid: false},
					Sort:     1,
					Active:   true,
					Display:  true,
				}, nil).Once()
				menuRepo.On("CountMenuChildren", testifymock.Anything, 1).Return(0, nil).Once()

				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, parentID).Return(model.Menu{
					ID:    parentID,
					Group: "settings",
				}, nil).Once()

				menuRepo.On("ReadNextSortForParentAndGroup", testifymock.Anything, int32(parentID), "settings").Return(7, nil).Once()

				menuRepo.
					On("UpdateMenuQuery", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(dto.UpdateMenuParams)
						assert.Equal(t, "settings", params.Group)
						assert.Equal(t, 7, params.Sort)
						assert.True(t, params.ParentID.Valid)
						assert.Equal(t, int32(parentID), params.ParentID.Int32)
					}).
					Return(nil).
					Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuRepo := mockrepo.NewMenuRepository()
			uc, mock, cleanup := newTestUseCaseWithMenuRepo(t, menuRepo)
			defer cleanup()

			tt.setup(mock, menuRepo)

			err := uc.UpdateMenuUseCase(context.Background(), tt.req, "user-1", tt.menuID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
