package menu

import (
	"context"
	"database/sql"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/menu/dto"
	menuPresenter "be-dashboard-nba/internal/presentation/menu/presenter"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestCreateMenuUseCase(t *testing.T) {
	parentID := 10

	tests := []struct {
		name    string
		req     menuPresenter.CreateMenuRequest
		setup   func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock)
		wantErr error
	}{
		{
			name: "positive/with-parent",
			req: menuPresenter.CreateMenuRequest{
				Name:     "Dashboard",
				Group:    "main",
				ParentID: &parentID,
				Active:   ptrBool(true),
				Display:  ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				menuRepo.On("ReadNextSortForParent", testifymock.Anything, int32(parentID)).Return(5, nil).Once()
				menuRepo.
					On("CreateMenuQuery", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(dto.CreateMenuParams)
						assert.Equal(t, 5, params.Sort)
						assert.True(t, params.ParentID.Valid)
						assert.Equal(t, int32(parentID), params.ParentID.Int32)
						assert.Equal(t, "main", params.Group)
					}).
					Return(nil).
					Once()
			},
		},
		{
			name: "edge/without-parent-uses-group-sort",
			req: menuPresenter.CreateMenuRequest{
				Name:    "Home",
				Group:   "main",
				Active:  ptrBool(true),
				Display: ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				menuRepo.On("ReadSortForGroup", testifymock.Anything, "main").Return(2, nil).Once()
				menuRepo.
					On("CreateMenuQuery", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(dto.CreateMenuParams)
						assert.Equal(t, 2, params.Sort)
						assert.False(t, params.ParentID.Valid)
						assert.Equal(t, "main", params.Group)
					}).
					Return(nil).
					Once()
			},
		},
		{
			name: "negative/sort-query-error",
			req: menuPresenter.CreateMenuRequest{
				Name:    "Home",
				Group:   "main",
				Active:  ptrBool(true),
				Display: ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuRepo.On("ReadSortForGroup", testifymock.Anything, "main").Return(0, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "negative/create-error",
			req: menuPresenter.CreateMenuRequest{
				Name:    "Home",
				Group:   "main",
				Active:  ptrBool(true),
				Display: ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuRepo.On("ReadSortForGroup", testifymock.Anything, "main").Return(2, nil).Once()
				menuRepo.On("CreateMenuQuery", testifymock.Anything, testifymock.Anything).Return(errors.New("insert failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "negative/begin-tx-error",
			req:  menuPresenter.CreateMenuRequest{Name: "X", Group: "main"},
			setup: func(m sqlmock.Sqlmock, _ *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin().WillReturnError(errors.New("begin failed"))
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "edge/rollback-error",
			req:  menuPresenter.CreateMenuRequest{Name: "X", Group: "main"},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback().WillReturnError(sql.ErrConnDone)

				menuRepo.On("ReadSortForGroup", testifymock.Anything, "main").Return(0, errors.New("sort failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuRepo := mockrepo.NewMenuRepository()
			uc, mock, cleanup := newTestUseCaseWithMenuRepo(t, menuRepo)
			defer cleanup()

			tt.setup(mock, menuRepo)

			err := uc.CreateMenuUseCase(context.Background(), tt.req, "user-1")
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func ptrBool(v bool) *bool { return &v }
