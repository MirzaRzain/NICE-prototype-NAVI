package container

import (
	"database/sql"

	"github.com/rs/zerolog"

	"be-dashboard-nba/internal/infrastructure/validator"
)

type Container struct {
	db        *sql.DB
	log       *zerolog.Logger
	validator *validator.Validator
}

func NewContainer(db *sql.DB, log *zerolog.Logger, validator *validator.Validator) *Container {
	return &Container{
		db:        db,
		log:       log,
		validator: validator,
	}
}

func (c *Container) GetDB() *sql.DB {
	return c.db
}

func (c *Container) GetLog() *zerolog.Logger {
	return c.log
}

func (c *Container) GetValidator() *validator.Validator {
	return c.validator
}
