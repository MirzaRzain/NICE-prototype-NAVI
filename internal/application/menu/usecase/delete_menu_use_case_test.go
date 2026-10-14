package menu

import (
	"context"
	"database/sql"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/domain/model"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestDeleteMenuUseCase(t *testing.T) {
	tests := []struct {
		name    string
		menuID  int
		setup   func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock)
		wantErr error
	}{
		{
			name:   "positive/delete-menu",
			menuID: 1,
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 1).Return(model.Menu{ID: 1}, nil).Once()
				menuRepo.On("DeleteMenuQuery", testifymock.Anything, 1).Return(nil).Once()
			},
		},
		{
			name:   "negative/not-found",
			menuID: 999,
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 999).Return(model.Menu{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrMenuIdNotFound,
		},
		{
			name:   "negative/delete-error",
			menuID: 1,
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 1).Return(model.Menu{ID: 1}, nil).Once()
				menuRepo.On("DeleteMenuQuery", testifymock.Anything, 1).Return(errors.New("delete failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name:   "edge/rollback-error",
			menuID: 1,
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback().WillReturnError(sql.ErrConnDone)

				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 1).Return(model.Menu{}, errors.New("read failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuRepo := mockrepo.NewMenuRepository()
			uc, mock, cleanup := newTestUseCaseForDeleteMenu(t, menuRepo)
			defer cleanup()

			tt.setup(mock, menuRepo)

			err := uc.DeleteMenuUseCase(context.Background(), tt.menuID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
