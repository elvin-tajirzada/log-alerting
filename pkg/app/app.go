package app

import (
	"fmt"
	"github.com/elvin-tacirzade/log-alerting/pkg/config"
	"github.com/elvin-tacirzade/log-alerting/pkg/contact-point"
	"github.com/elvin-tacirzade/log-alerting/pkg/db"
	"github.com/go-co-op/gocron"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type (
	App interface {
		Start()
		Shutdown()
	}

	app struct {
		Config    *config.Config
		Loki      *db.Loki
		Telegram  *contact_point.Telegram
		Scheduler *gocron.Scheduler
	}
)

func New() (App, error) {
	// create config
	conf, confErr := config.New()
	if confErr != nil {
		return nil, fmt.Errorf("failed to create a new config: %v", confErr)
	}

	// create loki
	loki := db.NewLoki(conf)

	// create telegram
	telegram, telegramErr := contact_point.NewTelegram(conf, loki)
	if telegramErr != nil {
		return nil, fmt.Errorf("failed to create a new telegram: %v", telegramErr)
	}

	scheduler := gocron.NewScheduler(conf.Timezone)

	return &app{
		Config:    conf,
		Loki:      loki,
		Telegram:  telegram,
		Scheduler: scheduler,
	}, nil
}

func (a *app) Start() {
	_, sendErrorMessageErr := a.Scheduler.Every(a.Config.EvaluationTime).Do(a.Telegram.SendErrorMessage, a.Config.EvaluationTime)
	if sendErrorMessageErr != nil {
		log.Fatalf("failed to setings up cron job for the SendErrorMessage func: %v", sendErrorMessageErr)
	}

	_, sendDailyReportErr := a.Scheduler.Every(1).Day().At("00:00").Do(a.Telegram.SendDailyReport, time.Hour*24)
	if sendDailyReportErr != nil {
		log.Fatalf("failed to settings up cron job for the SendDailyReport func: %v", sendDailyReportErr)
	}

	a.Scheduler.StartAsync()

	log.Println("log-alerting started successfully")
}

func (a *app) Shutdown() {
	signals := make(chan os.Signal)
	done := make(chan struct{})
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		a.Scheduler.Stop()
		if a.Scheduler.IsRunning() {
			log.Println("scheduler could not stopped")
		} else {
			log.Println("scheduler stopped")
		}

		log.Printf("received signal: %v\n", sig)
		log.Println("log-alerting exited. Bye...")
		done <- struct{}{}
	}()

	<-done
}
