package main

import (
	"log"
	"time"

	"be-dashboard-nba/internal/infrastructure/api"
	db "be-dashboard-nba/internal/infrastructure/database"
	"be-dashboard-nba/internal/infrastructure/logger"
	runtime "be-dashboard-nba/internal/infrastructure/runtime"
	"be-dashboard-nba/internal/infrastructure/runtime/container"
	"be-dashboard-nba/internal/infrastructure/validator"

	_ "be-dashboard-nba/docs"
)

// @title           Base APP Go
// @version         1.0
// @description     Template API for Golang project
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {your token}" into the field below.

func main() {
	if err := runtime.LoadConfig(); err != nil {
		log.Fatalf("ERROR loading config file: %s", err.Error())
	}

	setDefaultTimezone()

	logger.NewLogger()
	defer logger.Close()

	ctx, stop := runtime.NewRuntimeContext()
	defer stop()

	log := logger.WithContext(ctx)

	cfg := runtime.NewConfig()

	dbConn, err := db.NewDatabase(ctx)
	if err != nil {
		log.Fatal().Msgf("Failed to initialize database, %v", err)
	}

	validatorEngine := validator.NewValidator()
	c := container.NewContainer(dbConn, log.Raw(), validatorEngine)

	if err := api.RunFiberServer(ctx, c, cfg); err != nil {
		log.Fatal().Msgf("Failed to run server, %v", err)
	}
}

func setDefaultTimezone() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Now().Location()
	}

	time.Local = loc
}
