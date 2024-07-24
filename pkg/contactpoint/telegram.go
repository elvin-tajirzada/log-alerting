package contactpoint

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/elvin-tajirzada/log-alerting/internal/config"
	"github.com/elvin-tajirzada/log-alerting/internal/models"
	"github.com/elvin-tajirzada/log-alerting/pkg/db"
	"github.com/elvin-tajirzada/log-alerting/pkg/utils"
	"github.com/elvin-tajirzada/telegobot"
)

type Telegram struct {
	Telegobot       telegobot.Telegobot
	Loki            *db.Loki
	Timezone        *time.Location
	GrafanaPanelURL *url.URL
}

func NewTelegram(conf *config.Config, loki *db.Loki) (*Telegram, error) {
	t, err := telegobot.Start(conf.Telegram.Token, conf.Telegram.ChatID)
	if err != nil {
		return nil, fmt.Errorf("unable to start telegobot: %v", err)
	}

	return &Telegram{
		Telegobot:       t,
		Loki:            loki,
		Timezone:        conf.Timezone,
		GrafanaPanelURL: conf.GrafanaPanelURL,
	}, nil
}

func (t *Telegram) SendErrorMessage(evaluationTime time.Duration) {
	lokiLogEntry := t.getLokiData(evaluationTime)

	results := lokiLogEntry.Data.Result

	if len(results) > 0 {

		for _, result := range results {
			msg := t.createErrorMessage(result.Stream)

			t.send(msg)
		}

	}
}

func (t *Telegram) SendDailyReport(evaluationTime time.Duration) {
	lokiLogEntry := t.getLokiData(evaluationTime)

	if len(lokiLogEntry.Data.Result) == 0 {
		msg := `Alert Name: Daily Report 📅

No error 🔥

Keep going 👨🏻‍💻
`

		t.send(msg)
		return
	}

	groups := utils.GroupBy(lokiLogEntry)

	msg := t.createDailyReport(groups, len(lokiLogEntry.Data.Result))

	t.send(msg)
}

func (t *Telegram) getLokiData(evaluationTime time.Duration) *models.LokiLogEntry {
	lokiLogEntry, err := t.Loki.Get(evaluationTime)
	if err != nil {
		log.Fatalf("unable to get data from loki: %v", err)
	}

	return lokiLogEntry
}

func (t *Telegram) send(msg string) {

	if err := t.Telegobot.SendMessage(msg); err != nil {
		log.Fatalf("unable to send message to telegram: %v", err)
	}
}

func (t *Telegram) createErrorMessage(stream models.Stream) string {
	formatStr := `Alert Name: Server Error 💥

IP: %s

Device: %s

Method: %s

Status: %s

Path: %s

Message: %s

Time: %s

Panel URL: %s
`
	return fmt.Sprintf(formatStr,
		stream.IP,
		stream.Dt,
		stream.Method,
		stream.Status,
		stream.Path,
		stream.Msg,
		stream.Ts,
		t.GrafanaPanelURL.String(),
	)
}

func (t *Telegram) createDailyReport(groups map[string]map[string][]models.Stream, totalStreamCount int) string {
	formatStr := `Alert Name: Daily Report 📆

Total number of errors: ` + strconv.Itoa(totalStreamCount) + `

`
	for groupName, groupBy := range groups {
		formatStr += `Number of errors by ` + groupName + `:

`

		for name, streams := range groupBy {
			formatStr += `[` + name + `]: ` + strconv.Itoa(len(streams)) + `

`
		}
	}

	return formatStr
}
