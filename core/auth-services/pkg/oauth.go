package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/wafi11/workspaces/core/auth-services/config"
)

// --- GitHub ---

type githubUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func ExchangeGithubToken(ctx context.Context, conf config.SSOGithubConfig, code string) (string, error) {
	body, _ := json.Marshal(map[string]string{
		"client_id":     conf.ClientID,
		"client_secret": conf.ClientSecret,
		"code":          code,
	})

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access_token in github response")
	}
	return token, nil
}

func FetchGithubUser(ctx context.Context, token string) (*githubUser, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user githubUser
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}

// --- Google ---

type googleUser struct {
	ID    string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func ExchangeGoogleToken(ctx context.Context, conf config.SSOGoogleConfig, code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", conf.ClientID)
	data.Set("client_secret", conf.ClientSecret)
	data.Set("redirect_uri", conf.RedirectURL)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access_token in google response")
	}
	return token, nil
}

func FetchGoogleUser(ctx context.Context, token string) (*googleUser, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user googleUser
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}
