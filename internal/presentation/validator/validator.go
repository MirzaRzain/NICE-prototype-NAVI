package validator

import infra "be-dashboard-nba/internal/infrastructure/validator"

type Validator = infra.Validator

func New(engine *infra.Validator) *Validator {
	return engine
}
