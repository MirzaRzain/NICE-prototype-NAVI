package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"be-dashboard-nba/constant"
	masterDTO "be-dashboard-nba/internal/application/master/dto"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/validator"
	"be-dashboard-nba/internal/presentation/auth"
	"be-dashboard-nba/internal/presentation/master/presenter"
	requestPayload "be-dashboard-nba/internal/presentation/request"
	"be-dashboard-nba/internal/presentation/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ReadSources(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request presenter.ReadSourcesRequest
		if err := c.QueryParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Invalid query parameters"})
		}
		request.Init()
		if request.Page <= 0 {
			request.Page = constant.DefaultPage
		}
		if request.Limit <= 0 {
			request.Limit = constant.DefaultLimit
		}
		request.Offset = (request.Page - 1) * request.Limit

		sourceType := c.Query("source_type")
		if sourceType != "" && sourceType != "camera" && sourceType != "video" {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "source_type must be camera or video"})
		}
		startTime, err := parseSourceTime(c.Query("start_time"))
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "start_time must be RFC3339"})
		}
		endTime, err := parseSourceTime(c.Query("end_time"))
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "end_time must be RFC3339"})
		}
		if startTime != nil && endTime != nil && !startTime.Before(*endTime) {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "start_time must be before end_time"})
		}

		data, err := svc.ReadSources(c.UserContext(), masterDTO.ReadSourcesParams{SetSearch: request.SetSearch, Search: request.Search, SourceType: sourceType, StartTime: startTime, EndTime: endTime, Offset: request.Offset, Limit: request.Limit})
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get master sources"})
		}
		return response.Paginate(c, requestPayload.PaginationPayload{Page: request.Page, Limit: request.Limit}, int64(data.Pagination.TotalItems), response.DataPayload{Code: http.StatusOK, Data: data.Data, Message: "Successfully get master sources"})
	}
}

func parseSourceTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}
	t = t.UTC()
	return &t, nil
}

// CreateUsecase godoc
// @Summary      Create Master Usecase
// @Description  Creates the basic master usecase data. Inference-object assignments are configured separately.
// @Tags         Master Usecase
// @Accept       json
// @Produce      json
// @Param        request  body  presenter.CreateUsecaseRequest  true  "Usecase payload"
// @Success      201  {object}  response.DataPayload
// @Failure      400  {object}  response.ErrorPayload
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/usecases [post]
func CreateUsecase(svc contract.MasterUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request presenter.CreateUsecaseRequest
		if err := c.BodyParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Failed parse request"})
		}
		if err := validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}
		data, err := svc.CreateUsecase(c.UserContext(), request, "public")
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed create master usecase"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusCreated, Data: data, Message: "Successfully create master usecase"})
	}
}

// UpdateUsecase godoc
// @Summary      Update Master Usecase
// @Tags         Master Usecase
// @Accept       json
// @Produce      json
// @Param        usecase_id  path  int  true  "Usecase ID"
// @Param        request  body  presenter.UpdateUsecaseRequest  true  "Usecase payload"
// @Success      200  {object}  response.DataPayload
// @Failure      400  {object}  response.ErrorPayload
// @Failure      404  {object}  response.ErrorPayload
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/usecases/{usecase_id} [put]
func UpdateUsecase(svc contract.MasterUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		usecaseID, err := pathID(c, "usecase_id")
		if err != nil {
			return err
		}
		var request presenter.UpdateUsecaseRequest
		if err = c.BodyParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Failed parse request"})
		}
		if err = validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}
		if err = svc.UpdateUsecase(c.UserContext(), usecaseID, request, "public"); err != nil {
			if errors.Is(err, constant.ErrDataNotFound) {
				return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Master usecase not found"})
			}
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed update master usecase"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: nil, Message: "Successfully update master usecase"})
	}
}

// ReadUsecases godoc
// @Summary      Get Usecase List
// @Description  Returns the active usecases with their IDs, codes, and names.
// @Tags         Master Usecase
// @Produce      json
// @Success      200  {object}  presenter.UsecaseListResponse
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/usecases [get]
func ReadUsecases(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, err := svc.ReadUsecaseList(c.UserContext())
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get master usecases"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: data, Message: "Successfully get master usecases"})
	}
}

// ReadUsecasesWithConfig godoc
// @Summary      Get Usecases With Config
// @Description  Returns all active usecases and their key-value configurations.
// @Tags         Master Usecase
// @Produce      json
// @Success      200  {object}  presenter.UsecasesWithConfigResponse
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/usecases/config [get]
func ReadUsecasesWithConfig(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, err := svc.ReadUsecasesWithConfig(c.UserContext())
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get usecases with config"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: data, Message: "Successfully get usecases with config"})
	}
}

func ReadUsecase(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := pathID(c, "usecase_id")
		if err != nil {
			return err
		}
		data, err := svc.ReadUsecase(c.UserContext(), id)
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, constant.ErrDataNotFound) {
			return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Master usecase not found"})
		}
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get master usecase"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: data, Message: "Successfully get master usecase"})
	}
}

// ReadModels godoc
// @Summary      Get Master Models
// @Tags         Master Model
// @Produce      json
// @Success      200  {object}  response.DataPayload
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/models [get]
func ReadModels(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, err := svc.ReadModels(c.UserContext())
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get master models"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: data, Message: "Successfully get master models"})
	}
}

// ReadModelsByUsecaseCode godoc
// @Summary      Get Models By Usecase Code
// @Tags         Master Model
// @Produce      json
// @Param        usecase_code  path  string  true  "Usecase code"
// @Success      200  {object}  response.DataPayload
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/usecases/{usecase_code}/models [get]
func ReadModelsByUsecaseCode(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := strings.TrimSpace(c.Params("usecase_code"))
		if code == "" {
			return fiber.NewError(http.StatusBadRequest, "Invalid usecase_code")
		}
		data, err := svc.ReadModelsByUsecaseCode(c.UserContext(), code)
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get models by usecase"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: data, Message: "Successfully get models by usecase"})
	}
}

// ReadModel godoc
// @Summary      Get Master Model
// @Tags         Master Model
// @Produce      json
// @Param        model_id  path  int  true  "Model ID"
// @Success      200  {object}  response.DataPayload
// @Failure      404  {object}  response.ErrorPayload
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/models/{model_id} [get]
func ReadModel(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := pathID(c, "model_id")
		if err != nil {
			return err
		}
		data, err := svc.ReadModel(c.UserContext(), id)
		if errors.Is(err, sql.ErrNoRows) {
			return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Master model not found"})
		}
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get master model"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: data, Message: "Successfully get master model"})
	}
}

// CreateModel godoc
// @Summary      Create Master Model
// @Description  Upload a model and attach it to one or more global inference-object mappings.
// @Tags         Master Model
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Model file"
// @Param        name  formData  string  true  "Model display name"
// @Param        model_name  formData  string  true  "Model name"
// @Param        version  formData  string  true  "Model version"
// @Param        code  formData  string  true  "Model code"
// @Param        inference_object_ids  formData  string  true  "Comma-separated global inference-object IDs"
// @Success      201  {object}  response.DataPayload
// @Failure      400  {object}  response.ErrorPayload
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/models [post]
func CreateModel(svc contract.MasterUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Model file is required"})
		}
		request := presenter.CreateModelRequest{
			FileURL:   modelFileName(file.Filename),
			Name:      c.FormValue("name"),
			ModelName: c.FormValue("model_name"),
			Version:   c.FormValue("version"),
			Code:      c.FormValue("code"),
			Size:      file.Size,
		}
		request.InferenceObjectIDs, err = parseInferenceObjectIDs(c.FormValue("inference_object_ids"))
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Invalid inference_object_ids"})
		}
		if err := validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}

		storageDir := modelStorageDir()
		if err := os.MkdirAll(storageDir, 0750); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed prepare model storage"})
		}
		storedPath := filepath.Join(storageDir, request.FileURL)
		if err := c.SaveFile(file, storedPath); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed save model file"})
		}
		data, err := svc.CreateModel(c.UserContext(), request)
		if err != nil {
			_ = os.Remove(storedPath)
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed create master model"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusCreated, Data: data, Message: "Successfully create master model"})
	}
}

func DownloadModel(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := strings.TrimSpace(c.Params("usecase_code"))
		if code == "" {
			return fiber.NewError(http.StatusBadRequest, "Invalid usecase_code")
		}
		data, err := svc.ReadModelByUsecaseCode(c.UserContext(), code)
		if errors.Is(err, sql.ErrNoRows) {
			return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Model not found"})
		}
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get model"})
		}
		path := filepath.Join(modelStorageDir(), filepath.Base(data.FileURL))
		if _, err = os.Stat(path); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Model file not found"})
		}
		c.Attachment(data.Name + filepath.Ext(path))
		return c.SendFile(path, true)
	}
}

// DownloadModelByID godoc
// @Summary      Download Model By ID
// @Tags         Master Model
// @Produce      application/octet-stream
// @Param        model_id  path  int  true  "Model ID"
// @Success      200  {file}  binary
// @Failure      404  {object}  response.ErrorPayload
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/models/{model_id}/download [get]
func DownloadModelByID(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := pathID(c, "model_id")
		if err != nil {
			return err
		}
		data, err := svc.ReadModel(c.UserContext(), id)
		if errors.Is(err, sql.ErrNoRows) {
			return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Model not found"})
		}
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get model"})
		}
		path := filepath.Join(modelStorageDir(), filepath.Base(data.FileURL))
		if _, err = os.Stat(path); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Model file not found"})
		}
		c.Attachment(data.Name + filepath.Ext(path))
		return c.SendFile(path, true)
	}
}

// DownloadModelByInferenceObject godoc
// @Summary      Download Model By Inference Object
// @Tags         Master Model
// @Produce      application/octet-stream
// @Param        inference_object_id  path  int  true  "Usecase inference-object mapping ID"
// @Success      200  {file}  binary
// @Failure      404  {object}  response.ErrorPayload
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/models/download/inference-object/{inference_object_id} [get]
func DownloadModelByInferenceObject(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := pathID(c, "inference_object_id")
		if err != nil {
			return err
		}
		data, err := svc.ReadModelByInferenceObjectID(c.UserContext(), id)
		if errors.Is(err, sql.ErrNoRows) {
			return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Model not found"})
		}
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get model"})
		}
		path := filepath.Join(modelStorageDir(), filepath.Base(data.FileURL))
		if _, err = os.Stat(path); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Model file not found"})
		}
		c.Attachment(data.Name + filepath.Ext(path))
		return c.SendFile(path, true)
	}
}

func modelStorageDir() string {
	if dir := strings.TrimSpace(os.Getenv("MODEL_STORAGE_DIR")); dir != "" {
		return dir
	}
	return filepath.Join("storage", "models")
}

func modelFileName(original string) string {
	ext := filepath.Ext(filepath.Base(original))
	if len(ext) > 10 || strings.ContainsAny(ext, `/\\`) {
		ext = ""
	}
	return uuid.NewString() + strings.ToLower(ext)
}

// UpdateModel godoc
// @Summary      Update Master Model
// @Tags         Master Model
// @Accept       multipart/form-data
// @Produce      json
// @Param        model_id  path  int  true  "Model ID"
// @Param        file  formData  file  false  "Model file"
// @Param        name  formData  string  true  "Model display name"
// @Param        model_name  formData  string  true  "Model name"
// @Param        version  formData  string  true  "Model version"
// @Param        code  formData  string  true  "Model code"
// @Param        inference_object_ids  formData  string  true  "Comma-separated global inference-object IDs"
// @Success      200  {object}  response.DataPayload
// @Failure      400  {object}  response.ErrorPayload
// @Failure      404  {object}  response.ErrorPayload
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/models/{model_id} [put]
func UpdateModel(svc contract.MasterUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := pathID(c, "model_id")
		if err != nil {
			return err
		}
		var request presenter.UpdateModelRequest
		if err = c.BodyParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Failed parse request"})
		}
		if err = validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}
		if err = svc.UpdateModel(c.UserContext(), id, request); errors.Is(err, constant.ErrDataNotFound) {
			return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Master model not found"})
		} else if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed update master model"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: nil, Message: "Successfully update master model"})
	}
}

func DeleteModel(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := pathID(c, "model_id")
		if err != nil {
			return err
		}
		if err = svc.DeleteModel(c.UserContext(), id); errors.Is(err, constant.ErrDataNotFound) {
			return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Master model not found"})
		} else if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed delete master model"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: nil, Message: "Successfully delete master model"})
	}
}

// SetInferenceObjects godoc
// @Summary      Set Usecase Inference Objects
// @Description  Creates or replaces all inference-object assignments for a usecase.
// @Tags         Master Usecase
// @Accept       json
// @Produce      json
// @Param        usecase_id  path  int  true  "Usecase ID"
// @Param        request  body  presenter.SetInferenceObjectsRequest  true  "Inference-object assignments"
// @Success      200  {object}  response.DataPayload
// @Failure      400  {object}  response.ErrorPayload
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/usecases/{usecase_id}/inference-objects [put]
// @Router       /master/usecases/{usecase_id}/inference-objects [post]
func SetInferenceObjects(svc contract.MasterUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		usecaseID, err := pathID(c, "usecase_id")
		if err != nil {
			return err
		}
		var request presenter.SetInferenceObjectsRequest
		if err = c.BodyParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Failed parse request"})
		}
		if err = validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}
		ah, err := auth.GetAuth(c)
		if err != nil {
			return err
		}
		if err = svc.SetUsecaseInferenceObjects(c.UserContext(), usecaseID, request, ah.GetClaims().UserID); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed set inference objects"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: nil, Message: "Successfully set inference objects"})
	}
}

// ReadInferenceObjects godoc
// @Summary      Get Usecase Inference Objects
// @Tags         Master Usecase
// @Produce      json
// @Param        usecase_id  path  int  true  "Usecase ID"
// @Success      200  {object}  response.DataPayload
// @Failure      500  {object}  response.ErrorPayload
// @Router       /master/usecases/{usecase_id}/inference-objects [get]
func ReadInferenceObjects(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		usecaseID, err := pathID(c, "usecase_id")
		if err != nil {
			return err
		}
		data, err := svc.ReadUsecaseInferenceObjects(c.UserContext(), usecaseID)
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get inference objects"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: data, Message: "Successfully get inference objects"})
	}
}

func CreateVideo(svc contract.MasterUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request presenter.CreateVideoRequest
		if err := c.BodyParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Failed parse request"})
		}
		if err := validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}
		ah, err := auth.GetAuth(c)
		if err != nil {
			return err
		}
		data, err := svc.CreateVideo(c.UserContext(), request, ah.GetClaims().UserID)
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed create master video"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusCreated, Data: data, Message: "Successfully create master video"})
	}
}

func CreateCamera(svc contract.MasterUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request presenter.CreateCameraRequest
		if err := c.BodyParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Failed parse request"})
		}
		if err := validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}
		ah, err := auth.GetAuth(c)
		if err != nil {
			return err
		}
		data, err := svc.CreateCamera(c.UserContext(), request, ah.GetClaims().UserID)
		if err != nil {
			if errors.Is(err, constant.ErrFailedParseRequest) {
				return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Invalid installed_at"})
			}
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed create master camera"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusCreated, Data: data, Message: "Successfully create master camera"})
	}
}

func UpdateCamera(svc contract.MasterUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cameraID, err := pathID(c, "camera_id")
		if err != nil {
			return err
		}
		var request presenter.UpdateCameraRequest
		if err = c.BodyParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Failed parse request"})
		}
		if err = validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}
		ah, err := auth.GetAuth(c)
		if err != nil {
			return err
		}
		if err = svc.UpdateCamera(c.UserContext(), cameraID, request, ah.GetClaims().UserID); err != nil {
			if errors.Is(err, constant.ErrFailedParseRequest) {
				return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Invalid installed_at"})
			}
			if errors.Is(err, constant.ErrDataNotFound) {
				return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Master camera not found"})
			}
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed update master camera"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: nil, Message: "Successfully update master camera"})
	}
}

func ReadZones(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cameraID, err := pathID(c, "camera_id")
		if err != nil {
			return err
		}
		data, err := svc.ReadZonesByCamera(c.UserContext(), cameraID)
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed get camera zones"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: data, Message: "Successfully get camera zones"})
	}
}

func CreateZone(svc contract.MasterUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cameraID, err := pathID(c, "camera_id")
		if err != nil {
			return err
		}
		var request presenter.CreateZoneRequest
		if err = c.BodyParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Failed parse request"})
		}
		if err = validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}
		ah, err := auth.GetAuth(c)
		if err != nil {
			return err
		}
		data, err := svc.CreateZone(c.UserContext(), cameraID, request, ah.GetClaims().UserID)
		if err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed create camera zone"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusCreated, Data: data, Message: "Successfully create camera zone"})
	}
}

func UpdateZone(svc contract.MasterUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cameraID, err := pathID(c, "camera_id")
		if err != nil {
			return err
		}
		zoneID, err := pathID(c, "zone_id")
		if err != nil {
			return err
		}
		var request presenter.UpdateZoneRequest
		if err = c.BodyParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Failed parse request"})
		}
		if err = validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}
		ah, err := auth.GetAuth(c)
		if err != nil {
			return err
		}
		if err = svc.UpdateZone(c.UserContext(), cameraID, zoneID, request, ah.GetClaims().UserID); err != nil {
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed update camera zone"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: nil, Message: "Successfully update camera zone"})
	}
}

func DeleteZone(svc contract.MasterUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cameraID, err := pathID(c, "camera_id")
		if err != nil {
			return err
		}
		zoneID, err := pathID(c, "zone_id")
		if err != nil {
			return err
		}
		if err = svc.DeleteZone(c.UserContext(), cameraID, zoneID); err != nil {
			if errors.Is(err, constant.ErrDataNotFound) {
				return response.Error(c, response.ErrorParam{Code: http.StatusNotFound, Message: "Camera zone not found"})
			}
			return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "Failed delete camera zone"})
		}
		return response.Data(c, response.DataPayload{Code: http.StatusOK, Data: nil, Message: "Successfully delete camera zone"})
	}
}

func pathID(c *fiber.Ctx, name string) (int64, error) {
	id, err := strconv.ParseInt(c.Params(name), 10, 64)
	if err != nil || id < 1 {
		return 0, fiber.NewError(http.StatusBadRequest, "Invalid "+name)
	}
	return id, nil
}

func parseInferenceObjectIDs(value string) ([]int64, error) {
	parts := strings.Split(strings.TrimSpace(value), ",")
	ids := make([]int64, 0, len(parts))
	for _, part := range parts {
		id, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err != nil || id < 1 {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
