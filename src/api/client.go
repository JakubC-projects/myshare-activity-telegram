package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/JakubC-projects/myshare-activity-telegram/src/config"
	"github.com/JakubC-projects/myshare-activity-telegram/src/models"
	"github.com/goccy/go-json"
)

var baseUrl1 = config.Get().MyshareAPI.BaseUrl1
var baseUrl2 = config.Get().MyshareAPI.BaseUrl2

func get[T any](ctx context.Context, u models.User, url string) (T, error) {
	var res T

	token, err := getTokenSilently(ctx, u)
	if err != nil {
		return res, fmt.Errorf("cannot get token: %w", err)
	}

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

func post[T any](ctx context.Context, u models.User, url string, reqBody any) (T, error) {
	var res T

	token, err := getTokenSilently(ctx, u)
	if err != nil {
		return res, fmt.Errorf("cannot get token: %w", err)
	}

	var reqBodyReader io.Reader = nil
	if reqBody != nil {
		reqBodyData, err := json.Marshal(reqBody)
		if err != nil {
			return res, fmt.Errorf("cannot marshal request body: %w", err)
		}
		reqBodyReader = bytes.NewReader(reqBodyData)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBodyReader)
	if err != nil {
		return res, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, fmt.Errorf("request failed: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, fmt.Errorf("cannot read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(body))
		return res, fmt.Errorf("request failed: invalid status %d", resp.StatusCode)
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return res, fmt.Errorf("cannot unmarshal result: %w", err)
	}
	return res, nil
}
