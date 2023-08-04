package config

import (
	"fmt"
	"net/url"
	"os"
	"time"
)

const (
	lokiProtocol          = "http"
	lokiEndpoint          = "loki/api/v1/query_range"
	evaluationTimeDefault = time.Second * 10
)

type Config struct {
	EvaluationTime  time.Duration
	Timezone        *time.Location
	GrafanaPanelURL *url.URL
	Loki            loki
	Telegram        telegram
}

type telegram struct {
	Token  string
	ChatID string
}

type loki struct {
	URL string
}

func New() (*Config, error) {
	// evaluation time
	var (
		evaluationTime = os.Getenv("EVALUATION_TIME")
		evalTime       = evaluationTimeDefault
		evalTimeErr    error
	)

	if evaluationTime != "" {
		evalTime, evalTimeErr = time.ParseDuration(evaluationTime)
		if evalTimeErr != nil {
			return nil, fmt.Errorf("failed to parse EVALUATION_TIME: %v", evalTimeErr)
		}
	}

	// timezone
	timezone := os.Getenv("TIMEZONE")

	location, locationErr := time.LoadLocation(timezone)
	if locationErr != nil {
		return nil, fmt.Errorf("failed to load timezone: %v", locationErr)
	}

	// grafana panel URL
	grafanaPanelURL := os.Getenv("GRAFANA_PANEL_URL")
	u, uErr := url.ParseRequestURI(grafanaPanelURL)
	if uErr != nil {
		return nil, fmt.Errorf("failed to parse GRAFANA_PANEL_URL: %s", uErr)
	}

	// loki
	lokiHost := os.Getenv("LOKI_HOST")
	if lokiHost == "" {
		return nil, fmt.Errorf("environment parameter `LOKI_HOST` can't be empty")
	}

	lokiPort := os.Getenv("LOKI_PORT")
	if lokiPort == "" {
		return nil, fmt.Errorf("environment parameter `LOKI_PORT` can't be empty")
	}

	lokiURL := fmt.Sprintf("%s://%s:%s/%s", lokiProtocol, lokiHost, lokiPort, lokiEndpoint)

	// telegram
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		return nil, fmt.Errorf("environment parameter `TELEGRAM_TOKEN` can't be empty")
	}

	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	if telegramChatID == "" {
		return nil, fmt.Errorf("environment parameter `TELEGRAM_CHAT_ID` can't be empty")
	}

	return &Config{
		EvaluationTime:  evalTime,
		Timezone:        location,
		GrafanaPanelURL: u,
		Telegram: telegram{
			Token:  telegramToken,
			ChatID: telegramChatID,
		},
		Loki: loki{
			URL: lokiURL,
		},
	}, nil
}
