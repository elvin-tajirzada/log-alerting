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
	Loki            *Loki
	Telegram        *Telegram
	EvaluationTime  time.Duration
	Timezone        *time.Location
	GrafanaPanelURL *url.URL
}

type Telegram struct {
	Token  string
	ChatID string
}

type Loki struct {
	URL string
}

func New() (*Config, error) {
	evalTime, err := getEvaluationTime()
	if err != nil {
		return nil, err
	}

	timezone := os.Getenv("TIMEZONE")

	location, locationErr := time.LoadLocation(timezone)
	if locationErr != nil {
		return nil, fmt.Errorf("unable to load timezone: %v", locationErr)
	}

	grafanaPanelURL := os.Getenv("GRAFANA_PANEL_URL")
	u, uErr := url.ParseRequestURI(grafanaPanelURL)
	if uErr != nil {
		return nil, fmt.Errorf("unable to parse GRAFANA_PANEL_URL: %s", uErr)
	}

	lokiHost := os.Getenv("LOKI_HOST")
	if lokiHost == "" {
		return nil, fmt.Errorf("environment parameter `LOKI_HOST` can't be empty")
	}

	lokiPort := os.Getenv("LOKI_PORT")
	if lokiPort == "" {
		return nil, fmt.Errorf("environment parameter `LOKI_PORT` can't be empty")
	}

	lokiURL := fmt.Sprintf("%s://%s:%s/%s", lokiProtocol, lokiHost, lokiPort, lokiEndpoint)

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
		Telegram: &Telegram{
			Token:  telegramToken,
			ChatID: telegramChatID,
		},
		Loki: &Loki{
			URL: lokiURL,
		},
	}, nil
}

func getEvaluationTime() (time.Duration, error) {
	evaluationTime := os.Getenv("EVALUATION_TIME")
	if evaluationTime == "" {
		return evaluationTimeDefault, nil
	}

	evalTime, err := time.ParseDuration(evaluationTime)
	if err != nil {
		return 0, fmt.Errorf("unable to parse EVALUATION_TIME: %v", err)
	}

	return evalTime, nil
}
