package capture

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrBaseURLNotConfigured = errors.New("BASE_URL_CV is not configured")
	ErrInvalidPayload       = errors.New("invalid capture payload")
)

type Service struct {
	db         *sql.DB
	baseURL    string
	httpClient *http.Client
}

func NewService(db *sql.DB, baseURL string) *Service {
	return &Service{
		db:      db,
		baseURL: strings.TrimSpace(baseURL),
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (s *Service) Proxy(ctx context.Context, payload []byte) (int, string, []byte, error) {
	if s.baseURL == "" {
		return 0, "", nil, ErrBaseURLNotConfigured
	}

	enriched, err := s.enrichPayload(ctx, payload)
	if err != nil {
		return 0, "", nil, err
	}

	target, err := captureURL(s.baseURL)
	if err != nil {
		return 0, "", nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, target, strings.NewReader(string(enriched)))
	if err != nil {
		return 0, "", nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		return 0, "", nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, "", nil, err
	}
	contentType := response.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	return response.StatusCode, contentType, body, nil
}

func (s *Service) enrichPayload(ctx context.Context, payload []byte) ([]byte, error) {
	var root map[string]interface{}
	if err := json.Unmarshal(payload, &root); err != nil {
		return nil, ErrInvalidPayload
	}

	var camera map[string]interface{}
	if source, ok := root["source"].(map[string]interface{}); ok {
		camera, _ = source["camera"].(map[string]interface{})
	}

	if camera != nil {
		delete(camera, "port")
	}

	var identifier interface{}
	if camera != nil {
		identifier = camera["camera_id"]
	}
	if identifier == nil {
		identifier = root["camera_id"]
	}
	if identifier == nil || hasChannel(root, camera) {
		return json.Marshal(root)
	}

	channel, err := s.findChannel(ctx, identifier)
	if err != nil {
		return nil, err
	}
	if channel == "" {
		return json.Marshal(root)
	}
	if camera != nil {
		camera["channel"] = channel
	} else {
		root["channel"] = channel
	}
	return json.Marshal(root)
}

func hasChannel(root map[string]interface{}, camera map[string]interface{}) bool {
	if camera != nil {
		if value, exists := camera["channel"]; exists && value != nil {
			return true
		}
	}
	if value, exists := root["channel"]; exists && value != nil {
		return true
	}
	return false
}

func (s *Service) findChannel(ctx context.Context, identifier interface{}) (string, error) {
	if s.db == nil {
		return "", nil
	}
	value := fmt.Sprint(identifier)
	var channel sql.NullString
	var err error
	if id, parseErr := strconv.ParseInt(value, 10, 64); parseErr == nil {
		err = s.db.QueryRowContext(ctx, `SELECT channels FROM master_camera WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&channel)
	} else {
		err = s.db.QueryRowContext(ctx, `SELECT channels FROM master_camera WHERE camera_code=$1 AND deleted_at IS NULL`, value).Scan(&channel)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return channel.String, nil
}

func captureURL(base string) (string, error) {
	parsed, err := url.Parse(strings.TrimRight(base, "/"))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("invalid BASE_URL_CV")
	}
	path := strings.TrimRight(parsed.Path, "/")
	if strings.HasSuffix(path, "/api/v1") {
		parsed.Path = path + "/capture"
	} else {
		parsed.Path = path + "/api/v1/capture"
	}
	return parsed.String(), nil
}
