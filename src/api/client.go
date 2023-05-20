package api

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/goccy/go-json"
)

func get[T any](ctx context.Context, u models.User, path string) (T, error) {
	var res T

	token, err := getTokenSilently(ctx, u)
	if err != nil {
		return res, fmt.Errorf("cannot get token: %w", err)
	}

	url := config.Get().MyshareAPI.BaseUrl + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	if err != nil {
		return res, fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, fmt.Errorf("request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return res, fmt.Errorf("request failed: invalid status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, fmt.Errorf("cannot read body: %w", err)
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return res, fmt.Errorf("cannot unmarshal result: %w", err)
	}
	return res, nil
}
