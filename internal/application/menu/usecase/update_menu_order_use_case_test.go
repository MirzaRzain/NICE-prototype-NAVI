package menu

import (
	"context"
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

func TestUpdateMenuOrderUseCase(t *testing.T) {
	tests := []struct {
		name    string
		req     menuPresenter.UpdateMenuOrderRequest
		setup   func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock, updateSortCalls *int, updateChildCalls *int)
		wantErr error
		assert  func(menuRepo *mockrepo.MenuRepositoryMock, updateSortCalls int, updateChildCalls int)
	}{
		{
			name: "positive/update-multiple",
			req: menuPresenter.UpdateMenuOrderRequest{
				Group:     "main",
				SortedIDs: []int{3, 7},
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock, updateSortCalls *int, updateChildCalls *int) {
				m.ExpectBegin()
				m.ExpectCommit()

				menuRepo.
					On("UpdateMenuSortQuery", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						*updateSortCalls++
						params := args.Get(1).(dto.UpdateMenuSortParams)
						assert.Equal(t, "main", params.Group)
					}).
					Return(nil).
					Twice()
				menuRepo.
					On("UpdateChildrenGroup", testifymock.Anything, testifymock.Anything, "main").
					Run(func(args testifymock.Arguments) { *updateChildCalls++ }).
					Return(nil).
					Twice()
			},
			assert: func(_ *mockrepo.MenuRepositoryMock, updateSortCalls int, updateChildCalls int) {
				assert.Equal(t, 2, updateSortCalls)
				assert.Equal(t, 2, updateChildCalls)
			},
		},
		{
			name: "edge/empty-sorted-ids-commits-without-repo-calls",
			req: menuPresenter.UpdateMenuOrderRequest{
				Group:     "main",
				SortedIDs: []int{},
			},
			setup: func(m sqlmock.Sqlmock, _ *mockrepo.MenuRepositoryMock, _ *int, _ *int) {
				m.ExpectBegin()
				m.ExpectCommit()
			},
			assert: func(menuRepo *mockrepo.MenuRepositoryMock, updateSortCalls int, updateChildCalls int) {
				assert.Equal(t, 0, updateSortCalls)
				assert.Equal(t, 0, updateChildCalls)
				menuRepo.AssertNumberOfCalls(t, "UpdateMenuSortQuery", 0)
				menuRepo.AssertNumberOfCalls(t, "UpdateChildrenGroup", 0)
			},
		},
		{
			name: "negative/update-sort-error",
			req: menuPresenter.UpdateMenuOrderRequest{
				Group:     "main",
				SortedIDs: []int{3},
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock, _ *int, _ *int) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuRepo.On("UpdateMenuSortQuery", testifymock.Anything, testifymock.Anything).Return(errors.New("update sort failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuRepo := mockrepo.NewMenuRepository()
			uc, mock, cleanup := newTestUseCaseForUpdateMenuOrder(t, menuRepo)
			defer cleanup()

			updateSortCalls := 0
			updateChildCalls := 0
			tt.setup(mock, menuRepo, &updateSortCalls, &updateChildCalls)

			err := uc.UpdateMenuOrderUseCase(context.Background(), tt.req, "user-1")
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			if tt.assert != nil {
				tt.assert(menuRepo, updateSortCalls, updateChildCalls)
			}
		})
	}
}
