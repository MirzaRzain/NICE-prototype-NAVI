package menu

import (
	"context"
	"testing"

	"be-dashboard-nba/internal/domain/model"
	menuPresenter "be-dashboard-nba/internal/presentation/menu/presenter"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestReadListMenuUseCase(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(menuRepo *mockrepo.MenuRepositoryMock)
		wantLen int
		wantErr bool
	}{
		{
			name: "positive/returns-data",
			setup: func(menuRepo *mockrepo.MenuRepositoryMock) {
				menuRepo.On("ReadListMenuQuery", testifymock.Anything, testifymock.Anything).Return([]model.Menu{{ID: 1}}, nil).Once()
			},
			wantLen: 1,
		},
		{
			name: "edge/empty-list",
			setup: func(menuRepo *mockrepo.MenuRepositoryMock) {
				menuRepo.On("ReadListMenuQuery", testifymock.Anything, testifymock.Anything).Return([]model.Menu{}, nil).Once()
			},
			wantLen: 0,
		},
		{
			name: "negative/repo-error",
			setup: func(menuRepo *mockrepo.MenuRepositoryMock) {
				menuRepo.On("ReadListMenuQuery", testifymock.Anything, testifymock.Anything).Return(nil, errors.New("query failed")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuRepo := mockrepo.NewMenuRepository()
			uc, _, cleanup := newTestUseCaseForReadMenu(t, menuRepo)
			defer cleanup()

			tt.setup(menuRepo)

			data, err := uc.ReadListMenuUseCase(context.Background(), menuPresenter.ReadMenuListRequest{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, data, tt.wantLen)
		})
	}
}
